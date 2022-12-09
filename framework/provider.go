package framework

// NewInstance 定义如何创建一个新实例，所有服务容器的创建服务
type NewInstance func(...interface{}) (interface{}, error)

// ServiceProvider 定义一个服务提供者需要实现的接口
type ServiceProvider interface {
	// Register 在服务容器中注册一个实例化的方法，是否在注册的时候就实例化这个服务，需要参考IsDefer
	Register(Container) NewInstance
	// Boot 在调用实例化服务的时候会调用，可以把一些准备工作，如基础配置，初始化参数的操作放在这里
	// 如果Boot返回error，表示整个服务实例化失败
	Boot(Container) error
	// IsDefer 决定是否在注册的时候实例化这个服务，如果不注册的时候实例化，那就在第一次make的时候实例化
	// false 表示不需要延迟实例化
	IsDefer() bool
	// Params 定义传递给NewInstance的参数，可以自定义多个
	Params(Container) []interface{}
	// Name 代表这个服务提供者的凭证
	Name() string
}
