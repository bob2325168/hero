package middleware

import (
	"context"
	"fmt"
	"github.com/bob2325168/gohero/framework/gin"
	"log"
	"time"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	//使用函数回调
	return func(c *gin.Context) {
		finishCh := make(chan struct{}, 1)
		panicCh := make(chan interface{}, 1)

		// 执行业务逻辑前预操作，初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		// 一定跟着取消，尽早释放资源
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicCh <- p
				}
			}()
			//执行具体的业务逻辑
			c.Next()

			finishCh <- struct{}{}
		}()

		// 执行业务逻辑之后的操作
		select {
		case p := <-panicCh:
			log.Println(p)
			c.ISetStatus(500).IJson("time out")
		case <-finishCh:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.ISetStatus(500).IJson("time out")
		}
	}
}
