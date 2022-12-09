package contract

// AppKey 定义字符串凭证
const AppKey = "hero:app"

// App 定义接口
type App interface {
	// AppID 表示当前这个APP的唯一的ID，可以用于分布式锁
	AppID() string
	// Version 定义当前版本
	Version() string

	// BaseFolder 定义基础项目地址
	BaseFolder() string
	// ConfigFolder 定义配置文件地址
	ConfigFolder() string
	// LogFolder 定义日志文件地址
	LogFolder() string
	// ProviderFolder 定义业务自己的服务提供者地址
	ProviderFolder() string
	// MiddlewareFolder 定义业务定义的中间件
	MiddlewareFolder() string
	// CommandFolder 定义业务定义的命令地址
	CommandFolder() string
	// RuntimeFolder 定义业务的运行中间态信息
	RuntimeFolder() string
	// TestFolder 存放测试所需要的信息
	TestFolder() string

	// AppFolder 定义业务代码所在的目录，用于监控文件变更使用
	AppFolder() string
	// LoadAppConfig 加载新的AppConfig，key为对应函数转为小写下划线，ConfigFolder =>  config_folder
	LoadAppConfig(map[string]string)
}
