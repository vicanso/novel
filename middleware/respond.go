package middleware

import (
	"net/http"
	"strconv"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/util"
	"go.uber.org/zap"
)

// NewRespond 新建响应处理
func NewRespond() iris.Handler {
	return func(ctx iris.Context) {
		ctx.Next()
		logger := util.GetContextLogger(ctx)
		body := util.GetBody(ctx)
		if body == nil {
			return
		}
		// 对于>=400的错误记录出错数据
		if ctx.GetStatusCode() >= http.StatusBadRequest {
			logger.Error("request handle fail",
				zap.String("uri", ctx.Request().RequestURI),
				zap.Any("data", body),
			)
		}
		var err error
		contentType := ctx.GetContentType()
		switch body.(type) {
		case string:
			_, err = ctx.WriteString(body.(string))
		case []byte:
			if contentType == "" {
				ctx.ContentType(cs.ContentBinaryHeaderValue)
			}
			buf := body.([]byte)
			util.SetHeader(ctx, cs.HeaderContentLength, strconv.Itoa(len(buf)))
			_, err = ctx.Write(buf)
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
