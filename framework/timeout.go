package framework

import (
	"context"
	"fmt"
	"log"
	"time"
)

func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	//使用函数回调
	return func(c *Context) error {

		finishCh := make(chan struct{}, 1)
		panicCh := make(chan interface{}, 1)
		// 执行业务逻辑前预操作，初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.request.WithContext(durationCtx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicCh <- p
				}
			}()
			//执行具体的业务逻辑
			fun(c)

			finishCh <- struct{}{}
		}()

		// 执行业务逻辑之后的操作
		select {
		case p := <-panicCh:
			log.Println(p)
			c.responseWriter.WriteHeader(500)
		case <-finishCh:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.responseWriter.Write([]byte("time out"))
		}
		return nil
	}
}
