package trace

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
)

type HeroTraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a service instance
func (provider *HeroTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewHeroTraceService
}

// Boot will called when the service instantiate
func (provider *HeroTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *HeroTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *HeroTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

// / Name define the name for this service
func (provider *HeroTraceProvider) Name() string {
	return contract.TraceKey
}
