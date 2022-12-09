# Gin Default Server

This is API experiment for Gin.

```go
package main

import (
	"github.com/bob2325168/gohero/framework/gin"
	"github.com/bob2325168/gohero/framework/gin/ginS"
)

func main() {
	ginS.GET("/", func(c *gin.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
