package middleware

import (
	"time"

	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/xerror"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/vicanso/session"
)

const defaultMemoryStoreSize = 1024

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

// NewSession create a new session middleware
func NewSession(client *redis.Client, config SessionConfig) echo.MiddlewareFunc {
	var store session.Store
	if client != nil {
		store = session.NewRedisStore(client, nil)
	} else {
		store, _ = session.NewMemoryStore(defaultMemoryStoreSize)
	}
	opts := &session.Options{
		Store:        store,
		Key:          config.Cookie,
		MaxAge:       int(config.Expires.Seconds()),
		CookieKeys:   config.Keys,
		CookieMaxAge: int(config.CookieMaxAge.Seconds()),
		CookiePath:   config.CookiePath,
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if context.GetUserSession(c) != nil {
				return next(c)
			}
			res := c.Response()
			req := c.Request()
			sess := session.New(req, res, opts)
			_, err = sess.Fetch()
			if err != nil {
				err = xerror.NewSession(err.Error())
				return
			}
			us := service.NewUserSession(sess)
			context.SetUserSession(c, us)
			err = next(c)
			if err != nil {
				return
			}
			err = sess.Commit()
			if err != nil {
				err = xerror.NewSession(err.Error())
			}
			return
		}
	}
}
