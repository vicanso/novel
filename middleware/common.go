package middleware

import (
	"time"

	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/service"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/util"
)

var (
	// SessionHandler session处理
	SessionHandler iris.Handler
)

// 初始化session函数
func initSessionHandler() {
	client := service.GetRedisClient()
	defaultDuration := time.Hour * 24
	sessConfig := SessionConfig{
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
	SessionHandler = NewSession(client, sessConfig)
}

func init() {
	initSessionHandler()
}

// Session the session middleware
func Session(ctx iris.Context) {
	if SessionHandler == nil {
		resErr(ctx, util.ErrSessionIsNotReady)
		return
	}
	SessionHandler(ctx)
}

// IsLogined check login status，if not, will return error
func IsLogined(ctx iris.Context) {
	if util.GetAccount(ctx) == "" {
		resErr(ctx, util.ErrNeedLogined)
		return
	}
	ctx.Next()
}

// IsAnonymous check login status, if yes, will return error
func IsAnonymous(ctx iris.Context) {
	if util.GetAccount(ctx) != "" {
		resErr(ctx, util.ErrLoginedAlready)
		return
	}
	ctx.Next()
}

// WaitFor at least wait for duration
func WaitFor(d time.Duration) iris.Handler {
	ns := d.Nanoseconds()
	return func(ctx iris.Context) {
		start := time.Now()
		ctx.Next()
		use := time.Now().UnixNano() - start.UnixNano()
		if use < ns {
			time.Sleep(time.Duration(ns-use) * time.Nanosecond)
		}
	}
}

// IsSu check the user roles include su
func IsSu(ctx iris.Context) {
	account := util.GetAccount(ctx)
	if account == "" {
		resErr(ctx, util.ErrNeedLogined)
		return
	}
	sess := util.GetSession(ctx)
	roles := sess.GetStringSlice(cs.SessionRolesField)
	if !util.ContainsString(roles, model.UserRoleSu) {
		resErr(ctx, util.ErrUserForbidden)
		return
	}
	ctx.Next()
}

// IsNilQuery check the query is nil
func IsNilQuery(ctx iris.Context) {
	if ctx.Request().URL.RawQuery != "" {
		resErr(ctx, util.ErrQueryShouldBeNil)
		return
	}
	ctx.Next()
	return
}
