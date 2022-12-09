package env

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
)

type HeroEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *HeroEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewHeroEnv
}

// Boot will called when the service instantiate
func (provider *HeroEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *HeroEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *HeroEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

// Name define the name for this service
func (provider *HeroEnvProvider) Name() string {
	return contract.EnvKey
}
