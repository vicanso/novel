package main

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/config"
	_ "github.com/vicanso/novel/controller"
	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
	"go.uber.org/zap"
)

func main() {
	logger := util.GetLogger()
	redisClient := service.GetRedisClient()
	if redisClient != nil {
		// check the redis is healthy
		_, err := redisClient.Ping().Result()
		if err != nil {
			logger.Error("default redis client ping fail", zap.Error(err))
		}
	}

	app := iris.New()

	routes := router.List()

	app.Use(middleware.NewRecover())

	app.Use(middleware.NewRespond())

	app.Use(middleware.NewEntry())

	accessLogger := util.CreateAccessLogger()
	onStats := func(info *middleware.StatsInfo) {
		// TODO 可以写入至influxdb
		accessLogger.Info("",
			zap.String("trackId", info.TrackID),
			zap.String("account", info.Account),
			zap.String("ip", info.IP),
			zap.String("method", info.Method),
			zap.String("route", info.Route),
			zap.String("uri", info.URI),
			zap.Int("status", info.StatusCode),
			zap.Int("use", info.Consuming),
			zap.Int("type", info.Type),
			zap.Uint32("connecting", info.Connecting),
		)
		// accessLogger.Infof("%v", *info)
		// 如果觉得每次保存影响性能，可以只 % 10 == 0 才保存
		global.SaveConnectingCount(info.Connecting)
	}
	app.Use(middleware.NewStats(middleware.StatsConfig{
		OnStats: onStats,
	}))

	app.Use(middleware.NewLimiter(middleware.LimiterConfig{
		Max: 1000,
	}))

	app.Use(middleware.NewJSONParser(middleware.JSONParserConfig{}))

	// static file
	app.Get("/static/*", middleware.StaticServe(middleware.StaticServeConfig{
		Path: "./assets",
		// Asset: asset.New(),
		Header: map[string]string{
			"X-File": "My-Static-File",
		},
		Compression: true,
		MaxAge:      "24h",
		SMaxAge:     "1h",
	}))

	// method 不建议使用 any all
	routeInfos := make([]map[string]string, 0, 20)
	urlPrefix := config.GetString("urlPrefix")
	for i, r := range routes {
		// 对路由检测，判断是否有相同路由
		for j, tmp := range routes {
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
		m := map[string]string{
			"method": r.Method,
			"path":   r.Path,
		}
		routeInfos = append(routeInfos, m)
		routePath := urlPrefix + r.Path
		app.Handle(r.Method, routePath, r.Handlers...)
	}
	global.SaveRouteInfos(routeInfos)
	global.InitRouteCounter(routeInfos)

	global.StartApplication()
	app.Run(iris.Addr(config.GetString("listen")))
}
