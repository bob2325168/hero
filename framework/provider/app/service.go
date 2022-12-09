package app

import (
	"errors"
	"flag"
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/util"
	"github.com/google/uuid"
	"path/filepath"
)

// HeroApp 代表hero框架的app实现
type HeroApp struct {
	// 服务容器
	container framework.Container
	// 基础路径
	baseFolder string
	// appID
	appId string
	// 配置加载项
	configMap map[string]string
}

// NewHeroApp 初始化HeroApp
func NewHeroApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("params error")
	}
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	// 如果没有设置，则使用参数
	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
		flag.Parse()
	}
	appId := uuid.New().String()
	configMap := map[string]string{}
	return &HeroApp{
		container:  container,
		baseFolder: baseFolder,
		appId:      appId,
		configMap:  configMap,
	}, nil
}

// BaseFolder 基础目录
func (h *HeroApp) BaseFolder() string {

	if h.baseFolder != "" {
		return h.baseFolder
	}
	//如果参数也没有，默认当前路径
	return util.GetExecDirectory()
}

func (h *HeroApp) Version() string {
	return "0.0.1"
}

func (h *HeroApp) ConfigFolder() string {
	if val, ok := h.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (h *HeroApp) LogFolder() string {
	if val, ok := h.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "log")
}

func (h *HeroApp) HttpFolder() string {
	if val, ok := h.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "app", "http")
}

func (h *HeroApp) ConsoleFolder() string {
	if val, ok := h.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "app", "console")
}

func (h *HeroApp) StorageFolder() string {
	if val, ok := h.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "app", "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (h *HeroApp) ProviderFolder() string {
	if val, ok := h.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "app", "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (h *HeroApp) MiddlewareFolder() string {
	if val, ok := h.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(h.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (h *HeroApp) CommandFolder() string {
	if val, ok := h.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(h.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (h *HeroApp) RuntimeFolder() string {
	if val, ok := h.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(h.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (h *HeroApp) TestFolder() string {
	if val, ok := h.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "test")
}

// DeployFolder 定义测试需要的信息
func (h *HeroApp) DeployFolder() string {
	if val, ok := h.configMap["deploy_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "deploy")
}

// AppFolder 代表app目录
func (h *HeroApp) AppFolder() string {
	if val, ok := h.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "app")
}

// AppID 表示这个APP的唯一的ID
func (h *HeroApp) AppID() string {
	return h.appId
}

func (h *HeroApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		h.configMap[key] = val
	}
}
