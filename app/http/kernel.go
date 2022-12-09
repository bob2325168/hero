package http

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	// 设置为release模式，默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	// 默认启动一个Web引擎
	r := gin.New()
	// 设置engine
	r.SetContainer(container)
	// 使用recovery中间件
	r.Use(gin.Recovery())
	// 绑定路由
	Routes(r)
	return r, nil
}
