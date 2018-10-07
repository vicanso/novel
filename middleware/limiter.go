package middleware

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/vicanso/novel/global"
	"github.com/vicanso/novel/xerror"
	"github.com/vicanso/novel/xlog"

	"github.com/labstack/echo"
)

type (
	// LimiterConfig limiter config
	LimiterConfig struct {
		Max uint32
	}
)

var (
	errTooManyRequest = &xerror.HTTPError{
		StatusCode: http.StatusTooManyRequests,
		Message:    "too many request",
		Category:   xerror.ErrCategoryCommon,
	}
)

func resetApplicationStatus(d time.Duration) {
	// 等待后将程序重置为可用
	ticker := time.NewTicker(d)
	go func() {
		select {
		case <-ticker.C:
			// TODO 是否需要记录相关resume记录
			logger := xlog.Logger()
			logger.Info("application resume")
			global.StartApplication()
		}
	}()
}

// NewLimiter create a new limiter middleware
func NewLimiter(config LimiterConfig) echo.MiddlewareFunc {
	var count uint32
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				atomic.AddUint32(&count, ^uint32(0))
			}()
			v := atomic.AddUint32(&count, 1)
			if v > config.Max {
				err = errTooManyRequest
				// 如果多并发，还是会导致多个reset，影响不大，忽略
				if global.IsApplicationRunning() {
					global.PauseApplication()
					resetApplicationStatus(time.Second * 10)
				}
				return
			}
			return next(c)
		}
	}
}
