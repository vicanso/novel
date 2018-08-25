package middleware

import (
	"net/http"
	"time"

	"github.com/vicanso/novel/utils"

	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/vicanso/session"
)

type (
	// SessionConfig session's config
	SessionConfig struct {
		Cookie       string
		CookieMaxAge time.Duration
		CookiePath   string
		Expires      time.Duration
		Keys         []string
	}
)

// NewSession 创建新的session中间件
func NewSession(client *redis.Client, conf SessionConfig) iris.Handler {
	store := session.NewRedisStore(client, nil)
	opts := &session.Options{
		Store:        store,
		Key:          conf.Cookie,
		MaxAge:       int(conf.Expires.Seconds()),
		CookieKeys:   conf.Keys,
		CookieMaxAge: int(conf.CookieMaxAge.Seconds()),
		CookiePath:   conf.CookiePath,
	}
	return func(ctx iris.Context) {
		res := ctx.ResponseWriter()
		req := ctx.Request()
		sess := session.New(req, res, opts)
		_, err := sess.Fetch()
		if err != nil {
			utils.ResErr(ctx, &utils.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Category:   utils.ErrCategorySession,
				Message:    err.Error(),
				Code:       utils.ErrCodeSessionFetch,
			})
			return
		}
		utils.SetSession(ctx, sess)
		ctx.Next()
		err = sess.Commit()
		if err != nil {
			utils.ResErr(ctx, &utils.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Category:   utils.ErrCategorySession,
				Message:    err.Error(),
				Code:       utils.ErrCodeSessionCommit,
			})
			return
		}
	}
}
