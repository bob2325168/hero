package app

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
)

// HeroAppProvider APP的具体的实现方法
type HeroAppProvider struct {
	BaseFolder string
}

// Register 注册初始化服务实例方法
func (h *HeroAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewHeroApp
}

// Name 直接将服务对应的字符串凭证返回
func (h *HeroAppProvider) Name() string {
	return contract.AppKey
}

// Params 获取初始化参数
func (h *HeroAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

// IsDefer 是否延迟实例化，这里设置为true，表示将这个服务的实例化延迟到第一次make的时候
func (h *HeroAppProvider) IsDefer() bool {
	return false
}

// Boot 启动调用
func (h *HeroAppProvider) Boot(container framework.Container) error {
	return nil
}
