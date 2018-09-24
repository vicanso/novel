package router

import (
	"net/http"
	"strings"

	"github.com/vicanso/novel/global"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/util"
)

type (
	// Router 路由配置
	Router struct {
		Method   string
		Path     string
		Handlers []iris.Handler
	}
	// Group group router配置
	Group struct {
		Path     string
		Handlers []iris.Handler
	}
)

var (
	// routerList 路由列表
	routerList = make([]*Router, 0)
)

// Add 添加路由配置
func Add(method, path string, handlers ...iris.Handler) {
	r := &Router{
		Method:   strings.ToUpper(method),
		Path:     path,
		Handlers: handlers,
	}
	routerList = append(routerList, r)
}

// Add group add
func (g *Group) Add(method, path string, handlers ...iris.Handler) {
	currentPath := g.Path + path
	arr := make([]iris.Handler, len(g.Handlers))
	copy(arr, g.Handlers)
	arr = append(arr, handlers...)
	Add(method, currentPath, arr...)
}

// NewGroup 创建group
func NewGroup(path string, handlers ...iris.Handler) *Group {
	g := &Group{
		Path:     path,
		Handlers: handlers,
	}
	return g
}

// List 获取所有路由配置
func List() []*Router {
	return routerList
}

func init() {
	Add(http.MethodGet, "/ping", func(ctx iris.Context) {
		if global.IsApplicationRunning() {
			util.Res(ctx, "pong")
		} else {
			util.ResErr(ctx, util.ErrServiceUnavailable)
		}
	})
}
