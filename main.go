package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/asset"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/context"
	_ "github.com/vicanso/novel/controller"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/router"
	_ "github.com/vicanso/novel/schedule"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xlog"
	"go.uber.org/zap"
)

func healthCheck(listen string) {
	logger := xlog.Logger()
	url := "http://127.0.0.1" + listen + cs.PingRoute
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("health check fail",
			zap.Error(err),
		)
		os.Exit(1)
		return
	}
	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode >= 400 {
		logger.Error("helth check fail",
			zap.Int("status", statusCode),
		)
		os.Exit(1)
		return
	}
	os.Exit(0)
}

func main() {
	listen := config.GetString("listen")
	if util.ContainsString(os.Args, "--check=true") {
		healthCheck(listen)
		return
	}

	logger := xlog.Logger()
	influxdbClient := service.GetInfluxdbClient()
	if influxdbClient != nil {
		defer influxdbClient.Close()
	}
	// Echo instance
	e := echo.New()

	e.Use(middleware.NewRecover(middleware.RecoverConfig{}))

	e.Use(middleware.NewRespond(middleware.RespondConfig{}))

	e.Use(middleware.NewEntry(middleware.EntryConfig{}))

	accessLogger := xlog.AccessLogger()
	onStats := func(info *middleware.StatsInfo) {
		accessLogger.Info("",
			zap.String("trackId", info.TrackID),
			zap.String("requestId", info.RequestID),
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
		global.AddRouteCount(info.Method, info.Route)
		if influxdbClient != nil {
			tags := map[string]string{
				"method": info.Method,
				"route":  info.Route,
				"type":   strconv.Itoa(info.Type),
			}
			fields := map[string]interface{}{
				"trackId":    info.TrackID,
				"requestId":  info.RequestID,
				"account":    info.Account,
				"ip":         info.IP,
				"uri":        info.URI,
				"status":     info.StatusCode,
				"use":        info.Consuming,
				"connecting": info.Connecting,
			}
			err := service.WriteInfluxPoint(cs.MeasurementAccess, tags, fields)
			if err != nil {
				logger.Error("influxdb write point fail",
					zap.Error(err),
				)
			}
		}
	}
	e.Use(middleware.NewStats(middleware.StatsConfig{
		OnStats: onStats,
	}))

	e.Use(middleware.NewLimiter(middleware.LimiterConfig{
		Max: 1000,
	}))

	e.Use(middleware.NewJSONParser(middleware.JSONParserConfig{
		Limit: 100 * 1024,
	}))

	// TODO 是否需要增加ETag
	// 因为我使用的前置缓存Pike有ETag的处理，因此不需要添加

	e.GET("/", func(c echo.Context) (err error) {
		buf := asset.Get("index.html")
		context.SetContentType(c, "text/html")
		context.Res(c, buf)
		return
	})

	routes := router.List()
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
		routePath := r.Path
		if r.Path != cs.PingRoute {
			routePath = urlPrefix + r.Path
		}
		e.Add(r.Method, routePath, r.Handler, r.Mids...)
	}

	global.SaveRouteInfos(routeInfos)
	global.InitRouteCounter(routeInfos)

	global.StartApplication()

	defer logger.Sync()

	// Start server
	err := e.Start(listen)
	logger.Error("start server fail",
		zap.Error(err),
	)
}
