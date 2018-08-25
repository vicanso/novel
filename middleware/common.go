package middleware

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/utils"
)

// IsLogined check login status，if not, will return error
func IsLogined(ctx iris.Context) {
	if utils.GetAccount(ctx) == "" {
		utils.ResErr(ctx, utils.ErrNeedLogined)
		return
	}
	ctx.Next()
}
