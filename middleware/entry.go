package middleware

import (
	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
)

type (
	// EntryConfig 入口中间件的配置
	EntryConfig struct {
		Header map[string]string
	}
)

// NewEntry create a new entry middleware
func NewEntry(config EntryConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// 全局设置响应头
			for k, v := range config.Header {
				context.SetHeader(c, k, v)
			}
			// 所有的请求默认设置为no-cache
			context.SetNoCache(c)
			return next(c)
		}
	}
}
