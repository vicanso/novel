package middleware

import (
	"sync/atomic"
	"time"

	"github.com/kataras/iris"
)

type (
	// OnStats on stats function
	OnStats func(*StatsInfo)
	// StatsConfig stats的配置
	StatsConfig struct {
		OnStats OnStats
	}
	// StatsInfo 统计信息
	StatsInfo struct {
		IP         string
		TrackID    string
		Account    string
		Method     string
		Path       string
		URI        string
		StatusCode int
		Consuming  int
		Type       int
		Connecting uint32
	}
)

// NewStats 请求统计
func NewStats(conf StatsConfig) iris.Handler {
	var connectingCount uint32
	return func(ctx iris.Context) {
		atomic.AddUint32(&connectingCount, 1)
		startedAt := time.Now().UnixNano()
		ctx.Next()
		consuming := int(time.Now().UnixNano()-startedAt) / int(time.Millisecond)
		route := ctx.GetCurrentRoute()
		statusCode := ctx.GetStatusCode()

		info := &StatsInfo{
			URI:        ctx.Request().RequestURI,
			StatusCode: statusCode,
			Consuming:  consuming,
			Type:       statusCode / 100,
			Connecting: connectingCount,
			IP:         ctx.RemoteAddr(),
			TrackID:    getTrackID(ctx),
			Account:    getAccount(ctx),
		}
		if route != nil {
			info.Method = route.Method()
			info.Path = route.Path()
		}
		if conf.OnStats != nil {
			conf.OnStats(info)
		}
		atomic.AddUint32(&connectingCount, ^uint32(0))
	}
}
