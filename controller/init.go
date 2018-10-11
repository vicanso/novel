package controller

import (
	"regexp"
	"strconv"
	"time"

	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/xerror"
	"github.com/vicanso/novel/xlog"

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
	influxdbClient := service.GetInfluxdbClient()
	return middleware.NewTracker(middleware.TrackerConfig{
		OnTrack: func(info *middleware.TrackerInfo, c echo.Context) {

			us := context.GetUserSession(c)
			account := ""
			if us != nil {
				account = us.GetAccount()
			}

			postForm := make(map[string]interface{})
			json.Unmarshal(info.Form, &postForm)
			for k := range postForm {
				if defaultMaskFields.MatchString(k) {
					postForm[k] = "***"
				}
			}
			if len(postForm) == 0 {
				postForm = nil
			}
			err := ""
			if info.Err != nil {
				err = xerror.GetMessage(info.Err)
			}

			getContextLogger(c).Info("",
				zap.String("category", "tracker"),
				zap.String("account", account),
				zap.String("action", action),
				zap.Any("query", info.Query),
				zap.Any("params", info.Params),
				zap.Any("form", postForm),
				zap.Int("result", info.Result),
				zap.String("error", err),
			)
			if influxdbClient != nil {
				tags := map[string]string{
					"action": action,
					"result": strconv.Itoa(info.Result),
				}
				query, _ := json.Marshal(info.Query)
				params, _ := json.Marshal(info.Params)
				form, _ := json.Marshal(postForm)
				fields := map[string]interface{}{
					"trackId": context.GetTrackID(c),
				}

				if account != "" {
					fields["account"] = account
				}
				if len(query) != 0 {
					fields["query"] = string(query)
				}
				if len(params) != 0 {
					fields["params"] = string(params)
				}
				if len(form) != 0 {
					fields["form"] = string(form)
				}
				if err != "" {
					fields["error"] = err
				}
				err := service.WriteInfluxPoint(cs.MeasurementTracker, tags, fields)
				if err != nil {
					xlog.Logger().Error("influxdb write point fail",
						zap.Error(err),
					)
				}
			}
		},
	})
}
