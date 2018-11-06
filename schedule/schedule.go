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
	go initLatestCountResetTicker()
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
	logger := xlog.Logger()
	for range ticker.C {
		logger.Info(message + " schedule start")
		err := do()
		// TODO 检测不通过时，发送告警
		if err != nil {
			logger.Error(message+" fail",
				zap.Error(err),
			)
		}
		logger.Info(message + " schedule end")
	}
}

// 重置route count
func initRouteCountTicker() {
	// 每5分钟重置route count
	ticker := time.NewTicker(300 * time.Second)
	runTicker(ticker, "reset route count", func() error {
		global.ResetRouteCount()
		return nil
	}, initRedisCheckTicker)
}

// redis health check
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

// influxdb health check
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

// 定时更新book分类信息
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

// 定时更新章节信息
func initBookUpdateChaptersTicker() {
	interval := 1800 * time.Second
	ticker := time.NewTicker(interval)
	runTicker(ticker, "update book chapter", func() (err error) {
		b := service.Book{}
		status := strconv.Itoa(model.BookStatusPassed)
		books, err := b.List(&service.BookQueryParams{
			Status: status,
			Field:  "id,category,updatedAt",
		})
		if err != nil {
			return
		}
		now := time.Now().Unix()
		for _, book := range books {
			if util.ContainsString(book.Category, "完结") {
				continue
			}
			// 如果在上一次定时任务有更新，则此次跳过
			if now-book.UpdatedAt.Unix() < 2*int64(interval/time.Second) {
				continue
			}
			id := int(book.ID)
			key := fmt.Sprintf("%s-%d", cs.CacheBookUpdateChaptersLock, id)
			// 保证5分钟内只有实例获得更新
			ok, _ := service.Lock(key, 300*time.Second)
			if !ok {
				continue
			}
			b.UpdateChapters(id)
		}
		return
	}, initBookUpdateChaptersTicker)
}

// 定时重置最新用户行为count的信息
func initLatestCountResetTicker() {
	// 每小时重置一次
	ticker := time.NewTicker(3600 * time.Second)
	runTicker(ticker, "reset latest count", func() (err error) {
		// 保证30分钟内只有一个实例重置
		ok, _ := service.Lock(cs.CacheBookLatestCountRestLock, 1800*time.Second)
		if !ok {
			return
		}
		b := service.Book{}
		err = b.ResetLatestCount()
		return
	}, initLatestCountResetTicker)
}
