package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
)

type (
	// Router router
	Router struct {
		Method  string
		Path    string
		Mids    []echo.MiddlewareFunc
		Handler echo.HandlerFunc
	}
	// Group group router
	Group struct {
		Path string
		Mids []echo.MiddlewareFunc
	}
)

var (
	// routerList 路由列表
	routerList = make([]*Router, 0)
)

func init() {
	Add(http.MethodGet, "/ping", func(c echo.Context) (err error) {
		context.Res(c, "pong")
		return
	})
}

// Add add router config
func Add(method, path string, handler echo.HandlerFunc, mids ...echo.MiddlewareFunc) {
	r := &Router{
		Method:  strings.ToUpper(method),
		Path:    path,
		Mids:    mids,
		Handler: handler,
	}
	routerList = append(routerList, r)
}

// Add group add
func (g *Group) Add(method, path string, handler echo.HandlerFunc, mids ...echo.MiddlewareFunc) {
	currentPath := g.Path + path
	arr := make([]echo.MiddlewareFunc, len(g.Mids))
	copy(arr, g.Mids)
	arr = append(arr, mids...)
	Add(method, currentPath, handler, arr...)
}

// NewGroup create a group instance
func NewGroup(path string, mids ...echo.MiddlewareFunc) *Group {
	g := &Group{
		Path: path,
		Mids: mids,
	}
	return g
}

// List get all route list
func List() []*Router {
	return routerList
}
