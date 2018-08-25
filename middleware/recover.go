package middleware

import (
	"fmt"
	"net/http"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/utils"
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
			utils.SetNoCache(ctx)
			ctx.StatusCode(http.StatusInternalServerError)
			data := iris.Map{
				"message":   err.Error(),
				"exception": true,
			}
			if !utils.IsProduction() {
				data["stack"] = utils.GetStack(2 << 10)
			}
			ctx.JSON(data)
		}()
		ctx.Next()
	}
}
