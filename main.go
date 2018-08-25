package main

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/config"
	_ "github.com/vicanso/novel/controller"
	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/utils"
	"go.uber.org/zap"
)

func main() {
	logger := utils.GetLogger()
	app := iris.New()

	app.Use(middleware.NewRecover())

	app.Use(middleware.NewRespond())

	app.Use(middleware.NewEntry())

	accessLogger := utils.CreateAccessLogger()
	onStats := func(info *middleware.StatsInfo) {
		// TODO 可以写入至influxdb
		accessLogger.Infof("%v", *info)
	}
	app.Use(middleware.NewStats(middleware.StatsConfig{
		OnStats: onStats,
	}))

	app.Use(middleware.NewLimiter(middleware.LimiterConfig{
		Max: 1000,
	}))

	app.Use(middleware.NewJSONParser(middleware.JSONParserConfig{}))

	// method 不建议使用 any all
	for i, r := range router.List() {
		// 对路由检测，判断是否有相同路由
		for j, tmp := range router.List() {
			if j == i {
				continue
			}
			if r.Method == tmp.Method && r.Path == tmp.Path {
				logger.Error("duplicate route config",
					zap.String("method", r.Method),
					zap.String("path", r.Path),
				)
			}
		}
		app.Handle(r.Method, r.Path, r.Handlers...)
	}

	global.StartApplication()
	app.Run(iris.Addr(config.GetString("listen")))
}
