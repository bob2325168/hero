package kernel

import (
	"github.com/bob2325168/gohero/framework/gin"
	"net/http"
)

// HeroKernelService 引擎服务
type HeroKernelService struct {
	engine *gin.Engine
}

// NewHeroKernelService 初始化Web引擎服务实例
func NewHeroKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &HeroKernelService{engine: httpEngine}, nil
}

// HttpEngine 实现kernel入口的服务接口
func (s *HeroKernelService) HttpEngine() http.Handler {
	return s.engine
}
