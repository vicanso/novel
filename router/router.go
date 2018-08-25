package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/vicanso/novel/global"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/utils"
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
	// SessionHandler session处理
	SessionHandler iris.Handler
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

// 初始化session函数
func initSessionHandler() {
	client := service.GetRedisClient()
	defaultDuration := time.Hour * 24
	sessConfig := middleware.SessionConfig{
		// session cache expires
		Expires: config.GetDurationDefault("session.expires", defaultDuration),
		// the sesion cookie
		Cookie: config.GetSessionCookie(),
		// cookie max age (cookie有效期设置长一些)
		CookieMaxAge: 10 * config.GetDurationDefault("session.cookie.maxAge", defaultDuration),
		// cookie path
		CookiePath: config.GetCookiePath(),
		// cookie signed keys
		Keys: config.GetSessionKeys(),
	}
	SessionHandler = middleware.NewSession(client, sessConfig)
}

func init() {
	initSessionHandler()

	Add(http.MethodGet, "/ping", func(ctx iris.Context) {
		if global.IsApplicationRunning() {
			utils.Res(ctx, "pong")
		} else {
			utils.ResErr(ctx, utils.ErrServiceUnavailable)
		}
	})
}
