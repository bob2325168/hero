package framework

import (
	"log"
	"net/http"
	"strings"
)

// 框架核心结构
type Core struct {
	router      map[string]*Tree    //定义所有的routers
	middlewares []ControllerHandler // 从core这边设置的中间件
}

// 初始化核心架构
func NewCore() *Core {

	//将二级map写入一级map
	r := map[string]*Tree{}
	r["GET"] = NewTree()
	r["POST"] = NewTree()
	r["PUT"] = NewTree()
	r["DELETE"] = NewTree()

	return &Core{router: r}
}

// 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Core) Get(url string, handlers ...ControllerHandler) {
	// 将core的middleware和handlers结合
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 从core中初始化Group
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

//// 匹配路由，如果没有匹配返回nil
//func (c *Core) FindRouterByRequest(req *http.Request) []ControllerHandler {
//
//	//uri和method都转大写，保证大小写不敏感
//	uri := req.URL.Path
//	method := req.Method
//	upperMethod := strings.ToUpper(method)
//
//	//先找第一层map，然后再查找第二层
//	if methodHandlers, ok := c.router[upperMethod]; ok {
//		return methodHandlers.FindHandler(uri)
//	}
//	return nil
//}

func (c *Core) FindRouterNodeByRequest(req *http.Request) *node {

	//uri和method都转大写，保证大小写不敏感
	uri := req.URL.Path
	method := req.Method
	upperMethod := strings.ToUpper(method)

	//查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

// 框架核心架构handler
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// 封装自定义的context
	ctx := NewContext(request, response)

	//查找路由
	nd := c.FindRouterNodeByRequest(request)
	if nd == nil {
		ctx.SetStatus(404).Json("not found")
		return
	}

	// 设置context中的handlers字段
	ctx.SetHandlers(nd.handlers)

	//设置路由参数
	params := nd.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	//调用路由函数，如果返回err代表内部错误，返回500状态码
	if err := ctx.Next(); err != nil {
		ctx.SetStatus(500).Json("inner error")
		return
	}
}
