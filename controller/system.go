package controller

import (
	"os"
	"runtime"
	"time"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/asset"
	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/router"
)

var systemStartedAt = time.Now()

type (
	// SystemCtrl system controller
	SystemCtrl struct {
	}
	// StatusRes the response of status
	StatusRes struct {
		Status     string `json:"status,omitempty"`
		Uptime     string `json:"uptime,omitempty"`
		StartedAt  string `json:"startedAt,omitempty"`
		GoMaxProcs int    `json:"goMaxProcs,omitempty"`
		Version    string `json:"version,omitempty"`
		Pid        int    `json:"pid,omitempty"`
	}
	// StatsRes the response of stats
	StatsRes struct {
		Sys             uint64 `json:"sys,omitempty"`
		HeapSys         uint64 `json:"heapSys,omitempty"`
		HeapInuse       uint64 `json:"heapInuse,omitempty"`
		RoutineCount    int    `json:"routineCount,omitempty"`
		ConnectingCount uint32 `json:"connectingCount,omitempty"`
	}
	// RoutesRes the response of routes
	RoutesRes struct {
		Routes []map[string]string `json:"routes,omitempty"`
	}
)

func init() {
	ctrl := SystemCtrl{}
	system := router.NewGroup("/system")
	system.Add("GET", "/status", ctrl.getStatus)
	system.Add("GET", "/stats", ctrl.getStats)
	system.Add("GET", "/routes", ctrl.getRoutes)
	system.Add("GET", "/route-counts", ctrl.getRouteCounts)
	system.Add("GET", "/assets", ctrl.getAssets)
}

// getStatus get status info
func (sc *SystemCtrl) getStatus(c echo.Context) (err error) {
	status := "running"
	if !global.IsApplicationRunning() {
		status = "pause"
	}
	setCache(c, "10s")
	res(c, &StatusRes{
		Status:     status,
		Uptime:     time.Since(systemStartedAt).String(),
		StartedAt:  systemStartedAt.Format(time.RFC3339),
		GoMaxProcs: runtime.GOMAXPROCS(0),
		Version:    runtime.Version(),
		Pid:        os.Getpid(),
	})
	return
}

// getStats get stats info
func (sc *SystemCtrl) getStats(c echo.Context) (err error) {
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)
	var mb uint64 = 1024 * 1024
	setCache(c, "10s")
	res(c, &StatsRes{
		Sys:             mem.Sys / mb,
		HeapSys:         mem.HeapSys / mb,
		HeapInuse:       mem.HeapInuse / mb,
		RoutineCount:    runtime.NumGoroutine(),
		ConnectingCount: global.GetConnectingCount(),
	})
	return
}

// getRoutes get the route infos
func (sc *SystemCtrl) getRoutes(c echo.Context) (err error) {
	routeInfos := global.GetRouteInfos()
	setCache(c, "1m")
	res(c, &RoutesRes{
		Routes: routeInfos,
	})
	return
}

// getRouteCounts get route counts
func (sc *SystemCtrl) getRouteCounts(c echo.Context) (err error) {
	routeCountInfo := global.GetRouteCount()
	setCache(c, "1m")
	res(c, routeCountInfo)
	return
}

// getAssets get assets
func (sc *SystemCtrl) getAssets(c echo.Context) (err error) {
	adminAsset := asset.GetAdminAsset()
	webAsset := asset.GetWebAsset()
	setCache(c, "1m")
	res(c, map[string]interface{}{
		"admin": adminAsset.List(),
		"web": webAsset.List(),
	})
	return
}
