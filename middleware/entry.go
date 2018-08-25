package middleware

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/utils"
)

// NewEntry 创建新的entry
func NewEntry() iris.Handler {
	return func(ctx iris.Context) {
		utils.SetNoCache(ctx)
		logger := utils.CreateUserLogger(ctx)
		utils.SetContextLogger(ctx, logger)
		ctx.Next()
	}
}
