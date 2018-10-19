package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xerror"
	"go.uber.org/zap"
)

type (
	// RecoverConfig recover中间件的配置
	RecoverConfig struct {
	}
)

// NewRecover 创建新的recover中间件
func NewRecover(config RecoverConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				r := recover()
				if r == nil {
					return
				}
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				context.SetNoCache(c)
				message := err.Error()
				he := &xerror.HTTPError{
					StatusCode: http.StatusInternalServerError,
					Message:    message,
					Exception:  true,
					Category:   xerror.ErrCategoryPanic,
				}
				stack := util.GetStack(2, 7)
				if !util.IsProduction() {
					he.Stack = stack
				}
				logger := context.GetLogger(c)
				if logger != nil {
					logger.DPanic("exception error",
						zap.String("uri", c.Request().RequestURI),
						zap.Strings("stack", stack),
						zap.String("error", message),
					)
				}
				err = c.JSON(he.StatusCode, he)
				if err != nil {
					logger.Error("c.JSON fail",
						zap.Error(err),
					)
				}
			}()
			return next(c)
		}
	}
}
