package orm

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
)

type GormProvider struct {
}

func (gp *GormProvider) Register(container framework.Container) framework.NewInstance {
	return NewHeroGorm
}

func (gp *GormProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer ORM服务一定要延迟加载，不是必须的
func (gp *GormProvider) IsDefer() bool {
	return true
}

func (gp *GormProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (gp *GormProvider) Name() string {
	return contract.ORMKey
}
