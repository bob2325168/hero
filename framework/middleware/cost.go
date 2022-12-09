// Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package middleware

import (
	"github.com/bob2325168/gohero/framework/gin"
	"log"
	"time"
)

// Cost 记录调用链花费的时间
func Cost() gin.HandlerFunc {
	//使用函数回调
	return func(c *gin.Context) {
		//记录开始时间
		start := time.Now()
		log.Printf("api uri start: %v", c.Request.RequestURI)
		c.Next()

		//记录结束时间
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.Request.RequestURI, cost.Seconds())
	}
}
