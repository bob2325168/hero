package demo

import (
	"github.com/bob2325168/gohero/framework"
)

// DemoProvider 服务提供方
type DemoProvider struct {
	framework.ServiceProvider
	c framework.Container
}

// Name 直接将服务对应的字符串凭证返回
func (sp *DemoProvider) Name() string {
	return DemoKey
}

// Register 注册初始化服务实例方法
func (sp *DemoProvider) Register(c framework.Container) framework.NewInstance {
	return NewService
}

// IsDefer 是否延迟实例化，这里设置为true，表示将这个服务的实例化延迟到第一次make的时候
func (sp *DemoProvider) IsDefer() bool {
	return false
}

// Params 表示实例化的参数，这里只是实例化一个参数 container
func (sp *DemoProvider) Params(c framework.Container) []interface{} {
	return []interface{}{sp.c}
}

// Boot 只打印日志信息
func (sp *DemoProvider) Boot(c framework.Container) error {
	sp.c = c
	return nil
}
