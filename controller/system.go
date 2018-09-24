package controller

import (
	"os"
	"runtime"
	"time"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/router"
)

var systemStartedAt = time.Now()

type (
	// systemCtrl system controller
	systemCtrl struct {
	}
	statusRes struct {
		Status     string `json:"status,omitempty"`
		Uptime     string `json:"uptime,omitempty"`
		StartedAt  string `json:"startedAt,omitempty"`
		GoMaxProcs int    `json:"goMaxProcs,omitempty"`
		Version    string `json:"version,omitempty"`
		Pid        int    `json:"pid,omitempty"`
	}
	statsRes struct {
		Sys             uint64 `json:"sys,omitempty"`
		HeapSys         uint64 `json:"heapSys,omitempty"`
		HeapInuse       uint64 `json:"heapInuse,omitempty"`
		RoutineCount    int    `json:"routineCount,omitempty"`
		ConnectingCount uint32 `json:"connectingCount,omitempty"`
	}
	routesRes struct {
		Routes []map[string]string `json:"routes,omitempty"`
	}
)

func init() {
	ctrl := systemCtrl{}
	system := router.NewGroup("/system")
	system.Add("GET", "/status", ctrl.getStatus)
	system.Add("GET", "/stats", ctrl.getStats)
	system.Add("GET", "/routes", ctrl.getRoutes)
	system.Add("GET", "/route-counts", ctrl.getRouteCounts)
}

// getSystemStatus 获取系统状态信息
func (c *systemCtrl) getStatus(ctx iris.Context) {
	status := "running"
	if !global.IsApplicationRunning() {
		status = "pause"
	}
	setCache(ctx, "10s")
	res(ctx, &statusRes{
		Status:     status,
		Uptime:     time.Since(systemStartedAt).String(),
		StartedAt:  systemStartedAt.Format(time.RFC3339),
		GoMaxProcs: runtime.GOMAXPROCS(0),
		Version:    runtime.Version(),
		Pid:        os.Getpid(),
	})
}

// getSystemStats 获取系统性能信息
func (c *systemCtrl) getStats(ctx iris.Context) {
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)
	var mb uint64 = 1024 * 1024
	m := &statsRes{
		Sys:             mem.Sys / mb,
		HeapSys:         mem.HeapSys / mb,
		HeapInuse:       mem.HeapInuse / mb,
		RoutineCount:    runtime.NumGoroutine(),
		ConnectingCount: global.GetConnectingCount(),
	}
	setCache(ctx, "10s")
	res(ctx, m)
}

// getRoutes get the route infos
func (c *systemCtrl) getRoutes(ctx iris.Context) {
	routeInfos := global.GetRouteInfos()
	setCache(ctx, "1m")
	res(ctx, &routesRes{
		Routes: routeInfos,
	})
}

// getRouteCounts get route counts
func (c *systemCtrl) getRouteCounts(ctx iris.Context) {
	routeCountInfo := global.GetRouteCount()
	setCache(ctx, "1m")
	res(ctx, routeCountInfo)
}
