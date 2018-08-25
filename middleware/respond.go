package middleware

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/utils"
	"go.uber.org/zap"
)

// NewRespond 新建响应处理
func NewRespond() iris.Handler {
	logger := utils.GetLogger()
	return func(ctx iris.Context) {
		ctx.Next()
		body := utils.GetBody(ctx)
		if body == nil {
			return
		}
		var err error
		switch body.(type) {
		case string:
			_, err = ctx.WriteString(body.(string))
		case []byte:
			_, err = ctx.Binary(body.([]byte))
		default:
			_, err = ctx.JSON(body)
		}
		if err != nil {
			logger.Error("response fail",
				zap.String("uri", ctx.Request().RequestURI),
				zap.Error(err),
			)
		}
	}
}
