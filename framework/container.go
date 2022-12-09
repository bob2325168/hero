package framework

import (
	"errors"
	"fmt"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，进行替换并返回错误
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果关键字凭证未绑定一个服务提供者，就会panic
	// 所以在使用这个接口的时候要保证服务容器已经为这个关键字凭证绑定了一个服务提供者
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式，是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同的参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// HeroContainer 是服务容器的具体实现
type HeroContainer struct {
	// 强制要求HeroContainer实现Container接口
	Container
	// 存储注册的服务提供者，key为字符串凭证
	providers map[string]ServiceProvider
	// 存储具体的实例，key为字符串实例
	instances map[string]interface{}
	// 用于锁住对容器的变更操作，因为当Bind的时候会对实例或者服务提供者有一定的变动，需要使用一个机制来保证HeroContainer的并发性
	// 由于读多写少，使用读写锁优于互斥锁
	lock sync.RWMutex
}

// NewHeroContainer 创建一个服务容器
func NewHeroContainer() *HeroContainer {
	return &HeroContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// Bind 将服务容器和关键字进行绑定
func (hc *HeroContainer) Bind(provider ServiceProvider) error {

	hc.lock.Lock()
	// 获取服务容器的key
	key := provider.Name()
	hc.providers[key] = provider
	hc.lock.Unlock()

	// 如果provider不需要延迟实例化
	if provider.IsDefer() == false {
		if err := provider.Boot(hc); err != nil {
			return err
		}
		// 实例化方法
		params := provider.Params(hc)
		method := provider.Register(hc)
		instance, err := method(params...)
		if err != nil {
			fmt.Println("bind service provider ", key, "error: ", err)
			return errors.New(err.Error())
		}
		hc.instances[key] = instance
	}
	return nil
}

func (hc *HeroContainer) IsBind(key string) bool {
	return hc.findServiceProvider(key) != nil
}

func (hc *HeroContainer) findServiceProvider(key string) ServiceProvider {
	hc.lock.RLock()
	defer hc.lock.RUnlock()
	if sp, ok := hc.providers[key]; ok {
		return sp
	}
	return nil
}

// Make 调用内部实现的make方法
func (hc *HeroContainer) Make(key string) (interface{}, error) {
	return hc.make(key, nil, false)
}

func (hc *HeroContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return hc.make(key, params, true)
}

func (hc *HeroContainer) MustMake(key string) interface{} {
	serv, err := hc.make(key, nil, false)
	if err != nil {
		panic("container not contain key " + key)
	}
	return serv
}

// 真正的实例化一个服务
func (hc *HeroContainer) make(key string, params []interface{}, force bool) (interface{}, error) {

	hc.lock.RLock()
	defer hc.lock.RUnlock()
	// 查询是否已经注册了这个服务提供者，如果没有注册，返回错误
	sp := hc.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}
	// 强制实例化，需要重新实例化新的实例
	if force {
		return hc.newInstance(sp, params)
	}
	// 不需要强制实例化，那么在容器中已经实例化了就直接取出来使用即可
	if ins, ok := hc.instances[key]; ok {
		return ins, nil
	}
	// 如果容器中还没有实例化，就进行实例化
	inst, err := hc.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	hc.instances[key] = inst
	return inst, nil
}

func (hc *HeroContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {

	if err := sp.Boot(hc); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(hc)
	}
	method := sp.Register(hc)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

// NameList 列出容器中所有服务提供者的字符串凭证
func (hc *HeroContainer) NameList() []string {
	var ret []string
	for _, provider := range hc.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}
