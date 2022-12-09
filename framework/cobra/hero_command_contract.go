package cobra

import "github.com/bob2325168/gohero/framework/contract"

// MustMakeApp 从容器中获取App服务
func (c *Command) MustMakeApp() contract.App {
	return c.GetContainer().MustMake(contract.AppKey).(contract.App)
}

// MustMakeKernel 从容器中获取Kernel服务
func (c *Command) MustMakeKernel() contract.Kernel {
	return c.GetContainer().MustMake(contract.KernelKey).(contract.Kernel)
}
