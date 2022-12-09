package contract

import "net/http"

// KernelKey 提供kernel服务凭证
const KernelKey = "hero:kernel"

// Kernel 提供框架最核心的结构
type Kernel interface {
	// HttpEngine 作为net/http框架使用，实际上是gin.Engine
	HttpEngine() http.Handler
}
