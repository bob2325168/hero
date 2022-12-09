package id

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
)

type HeroIDProvider struct {
}

// Register registe a new function for make a service instance
func (provider *HeroIDProvider) Register(c framework.Container) framework.NewInstance {
	return NewHeroIDService
}

// Boot will called when the service instantiate
func (provider *HeroIDProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *HeroIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *HeroIDProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

// Name define the name for this service
func (provider *HeroIDProvider) Name() string {
	return contract.IDKey
}
