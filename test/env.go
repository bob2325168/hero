package tests

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/provider/app"
	"github.com/bob2325168/gohero/framework/provider/env"
)

const (
	BasePath = "/Users/yejianfeng/Documents/UGit/coredemo/"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewHeroContainer()
	// 绑定App服务提供者
	container.Bind(&app.HeroAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.HeroEnvProvider{})
	return container
}
