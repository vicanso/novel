package controller

import (
	"regexp"
	"time"

	"github.com/vicanso/novel/xerror"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
)

var (
	json                = jsoniter.ConfigCompatibleWithStandardLibrary
	getRequestQuery     = context.GetRequestQuery
	getRequestBody      = context.GetRequestBody
	res                 = context.Res
	resCreated          = context.ResCreated
	setCache            = context.SetCache
	setCacheWithSMaxAge = context.SetCacheWithSMaxAge
	setHeader           = context.SetHeader
	setPrivateCache     = context.SetPrivateCache
	setContentType      = context.SetContentType
	getUserSession      = context.GetUserSession
	now                 = util.Now
	getTrackID          = context.GetTrackID
	userSession         = initSessionHandler()
	userService         = &service.User{}
	bookService         = &service.Book{}
	getContextLogger    = context.GetLogger
)

// 初始化session函数
func initSessionHandler() echo.MiddlewareFunc {
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
	return middleware.NewSession(client, sessConfig)
}

// createTracker create a tracker middleware
func createTracker(action string) echo.MiddlewareFunc {
	defaultMaskFields := regexp.MustCompile(`password`)
	return middleware.NewTracker(middleware.TrackerConfig{
		OnTrack: func(info *middleware.TrackerInfo, c echo.Context) {
			form := make(map[string]interface{})
			json.Unmarshal(info.Form, &form)
			for k := range form {
				if defaultMaskFields.MatchString(k) {
					form[k] = "***"
				}
			}
			if len(form) == 0 {
				form = nil
			}
			err := ""
			if info.Err != nil {
				err = xerror.GetMessage(info.Err)
			}

			getContextLogger(c).Info("",
				zap.String("category", "tracker"),
				zap.String("action", action),
				zap.Any("query", info.Query),
				zap.Any("params", info.Params),
				zap.Any("form", form),
				zap.Int("result", info.Result),
				zap.String("error", err),
			)
		},
	})
}
