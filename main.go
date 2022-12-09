// Copyright 2022 bob2325168.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/bob2325168/gohero/app/console"
	"github.com/bob2325168/gohero/app/http"
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/provider/app"
	"github.com/bob2325168/gohero/framework/provider/cache"
	"github.com/bob2325168/gohero/framework/provider/config"
	"github.com/bob2325168/gohero/framework/provider/distributed"
	"github.com/bob2325168/gohero/framework/provider/env"
	"github.com/bob2325168/gohero/framework/provider/id"
	"github.com/bob2325168/gohero/framework/provider/kernel"
	"github.com/bob2325168/gohero/framework/provider/log"
	"github.com/bob2325168/gohero/framework/provider/orm"
	"github.com/bob2325168/gohero/framework/provider/redis"
	"github.com/bob2325168/gohero/framework/provider/trace"
)

func main() {
	// 初始化服务容器
	container := framework.NewHeroContainer()
	// 绑定APP服务提供者
	container.Bind(&app.HeroAppProvider{})
	// 后续初始化需要绑定的服务提供者
	container.Bind(&env.HeroEnvProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&config.HeroConfigProvider{})
	container.Bind(&id.HeroIDProvider{})
	container.Bind(&trace.HeroTraceProvider{})
	container.Bind(&log.HeroLogServiceProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&redis.RedisProvider{})
	container.Bind(&cache.HeroCacheProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.HeroKernelProvider{HttpEngine: engine})
	}

	// 创建和运行根command
	console.RunCommand(container)
}
