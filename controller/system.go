package controller

import (
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
)

func init() {
	ctrl := systemCtrl{}
	system := router.NewGroup("/system")
	system.Add("GET", "/status", ctrl.getStatus)
	system.Add("GET", "/stats", ctrl.getStats)
}

// getSystemStatus 获取系统状态信息
func (c *systemCtrl) getStatus(ctx iris.Context) {
	status := "running"
	if !global.IsApplicationRunning() {
		status = "pause"
	}
	m := iris.Map{
		"status":     status,
		"uptime":     time.Since(systemStartedAt).String(),
		"startedAt":  systemStartedAt.Format(time.RFC3339),
		"goMaxProcs": runtime.GOMAXPROCS(0),
		"version":    runtime.Version(),
	}
	setCache(ctx, "10s")
	res(ctx, m)
}

// getSystemStats 获取系统性能信息
func (c *systemCtrl) getStats(ctx iris.Context) {
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)
	var mb uint64 = 1024 * 1024
	m := iris.Map{
		"sys":          mem.Sys / mb,
		"heapSys":      mem.HeapSys / mb,
		"heapInuse":    mem.HeapInuse / mb,
		"routineCount": runtime.NumGoroutine(),
	}
	setCache(ctx, "10s")
	res(ctx, m)
}
