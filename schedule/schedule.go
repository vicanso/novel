package service

import (
	"time"

	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/xlog"
	"go.uber.org/zap"
)

func init() {
	go initRouteCountTicker()
	go initRedisCheckTicker()
	go initInfluxdbCheckTicker()
}

func runTicker(ticker *time.Ticker, message string, do func() error, restart func()) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			xlog.Logger().DPanic(message+" panic",
				zap.Error(err),
			)
		}
		// 如果退出了，重新启动
		go restart()
	}()
	for range ticker.C {
		err := do()
		// TODO 检测不通过时，发送告警
		if err != nil {
			xlog.Logger().Error(message+" fail",
				zap.Error(err),
			)
		}
	}
}

func initRouteCountTicker() {
	// 每5分钟重置route count
	ticker := time.NewTicker(300 * time.Second)
	runTicker(ticker, "reset route count", func() error {
		global.ResetRouteCount()
		return nil
	}, initRedisCheckTicker)
}

func initRedisCheckTicker() {
	client := service.GetRedisClient()
	// 未使用redis，则不需要检测
	if client == nil {
		return
	}
	// 每一分钟检测一次
	ticker := time.NewTicker(6 * time.Second)
	runTicker(ticker, "redis check", func() error {
		_, err := client.Ping().Result()
		return err
	}, initRedisCheckTicker)
}

func initInfluxdbCheckTicker() {
	clinet := service.GetInfluxdbClient()
	if clinet == nil {
		return
	}

	// 每一分钟检测一次
	ticker := time.NewTicker(60 * time.Second)
	runTicker(ticker, "influxdb check", func() error {
		_, _, err := clinet.Ping(3 * time.Second)
		return err
	}, initInfluxdbCheckTicker)
}
