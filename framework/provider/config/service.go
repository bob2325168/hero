package config

import (
	"bytes"
	"fmt"
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type HeroConfig struct {
	c framework.Container
	// 文件夹
	folder string
	// 路径的分隔符，默认为点
	keyDelim string
	//配置文件锁
	lock sync.RWMutex
	// 所有环境变量
	envMaps map[string]string
	// 配置文件结构，key为文件名
	confMaps map[string]interface{}
	// 配置文件的原始信息
	confRaws map[string][]byte
}

// 读取某个文件配置
func (conf *HeroConfig) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	//  判断文件是否以yaml或者yml作为后缀
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		// 读取文件内容
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		// 直接针对文本做环境变量的替换
		bf = replace(bf, conf.envMaps)
		// 解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		conf.confMaps[name] = c
		conf.confRaws[name] = bf

		// 读取app.path中的信息，更新app对应的folder
		if name == "app" && conf.c.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appService := conf.c.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}

// 删除文件操作
func (conf *HeroConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

func NewHeroConfig(params ...interface{}) (interface{}, error) {

	c := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	//检查文件夹是否存在
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
	}

	heroConf := &HeroConfig{
		c:        c,
		folder:   envFolder,
		envMaps:  envMaps,
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
		keyDelim: ".",
		lock:     sync.RWMutex{},
	}

	// 读取每个文件
	files, err := os.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range files {
		fileName := file.Name()
		err := heroConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// 监控文件夹文件是否变化
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watch.Add(envFolder)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				{
					//判断事件发生类型，create创建，write写入，remove删除
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						heroConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						heroConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						heroConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error: ", err)
					return
				}
			}
		}
	}()
	return heroConf, nil
}

// 表示使用环境变量替换context中的env(XXX)的环境变量
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}
	// 直接使用ReplaceAll替换，虽然性能可能不太好，但是频率低所以可以接受
	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

// 查找某个路径的配置项
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// 判断这个路径是否为1
		if len(path) == 1 {
			return next
		}

		// 判断下一个路径的类型
		switch next.(type) {
		case map[interface{}]interface{}:
			// 如果是interface类型的map，使用cast进行value转换
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// 如果是map[string]直接循环调用
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil
		}
	}
	return nil
}

// 通过path获取某个元素
func (conf *HeroConfig) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

// IsExist check setting is exist
func (conf *HeroConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

// Get a new interface
func (conf *HeroConfig) Get(key string) interface{} {
	return conf.find(key)
}

// GetBool get bool type
func (conf *HeroConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

// GetInt get Int type
func (conf *HeroConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *HeroConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *HeroConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *HeroConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *HeroConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *HeroConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *HeroConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *HeroConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *HeroConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

// Load a config to a struct, val should be an pointer
func (conf *HeroConfig) Load(key string, val interface{}) error {

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(conf.find(key))
}
