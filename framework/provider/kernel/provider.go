package kernel

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
	"github.com/bob2325168/gohero/framework/gin"
)

// HeroKernelProvider 提供web引擎
type HeroKernelProvider struct {
	HttpEngine *gin.Engine
}

// Register 注册服务提供者
func (p *HeroKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewHeroKernelService
}

// Boot 启动时判断是否由外界注入了Engine，如果有注入使用注入的，如果没有重新实例化
func (p *HeroKernelProvider) Boot(c framework.Container) error {
	if p.HttpEngine == nil {
		p.HttpEngine = gin.Default()
	}
	p.HttpEngine.SetContainer(c)
	return nil
}

func (p *HeroKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{p.HttpEngine}
}

func (p *HeroKernelProvider) IsDefer() bool {
	return false
}

func (p *HeroKernelProvider) Name() string {
	return contract.KernelKey
}
