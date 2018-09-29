package middleware

import (
	"fmt"
	"net/http"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/util"
	"go.uber.org/zap"
)

// NewRecover 创建异常恢复中间件
func NewRecover() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			if ctx.IsStopped() {
				return
			}
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			util.SetNoCache(ctx)
			ctx.StatusCode(http.StatusInternalServerError)
			data := iris.Map{
				"message":   err.Error(),
				"exception": true,
			}
			stack := util.GetStack(2, 7)
			if !util.IsProduction() {
				data["stack"] = stack
			}
			ctx.JSON(data)
			util.GetContextLogger(ctx).
				Error("exception error",
					zap.String("uri", ctx.Request().RequestURI),
					zap.Strings("stack", stack),
					zap.String("error", err.Error()),
				)
		}()
		ctx.Next()
	}
}
