package service

import (
	"time"

	"github.com/vicanso/novel/global"
)

func init() {
	go initRouteCountTicker()
}

func initRouteCountTicker() {
	// 每5分钟重置route count
	ticker := time.NewTicker(300 * time.Second)
	for range ticker.C {
		global.ResetRouteCount()
	}
}
