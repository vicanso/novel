package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xlog"
	"go.uber.org/zap"
)

func init() {
	if util.IsDevelopment() {
		return
	}
	go initRouteCountTicker()
	go initRedisCheckTicker()
	go initInfluxdbCheckTicker()
	go initBookCategoryTicker()
	go initBookUpdateChaptersTicker()
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

func initBookCategoryTicker() {
	ticker := time.NewTicker(600 * time.Second)
	runTicker(ticker, "update book category", func() error {
		// 避免启动多个实例时并发更新
		ok, _ := service.Lock(cs.CacheBookCategoriesLock, 300*time.Second)
		if !ok {
			return nil
		}
		b := service.Book{}
		return b.UpdateCategories()
	}, initBookCategoryTicker)
}

func initBookUpdateChaptersTicker() {
	ticker := time.NewTicker(3600 * time.Second)
	runTicker(ticker, "update book chapter", func() (err error) {
		b := service.Book{}
		status := strconv.Itoa(model.BookStatusPassed)
		books, err := b.List(&service.BookQueryParams{
			Status: status,
			Field:  "id",
		})
		if err != nil {
			return
		}
		for _, book := range books {
			id := int(book.ID)
			key := fmt.Sprintf("%s-%d", cs.CacheBookUpdateChaptersLock, id)
			// 保证30分钟内只有实例获得更新
			ok, _ := service.Lock(key, 1800*time.Second)
			if !ok {
				return
			}
			b.UpdateChapters(id, 10)
		}
		return
	}, initBookUpdateChaptersTicker)
}
