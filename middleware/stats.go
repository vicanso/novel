package middleware

import (
	"net/url"
	"sync/atomic"
	"time"

	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/xerror"

	"github.com/labstack/echo"
)

type (
	// OnStats on stats function
	OnStats func(*StatsInfo)
	// StatsConfig stats config
	StatsConfig struct {
		OnStats OnStats
	}
	// StatsInfo 统计信息
	StatsInfo struct {
		RequestID  string
		IP         string
		TrackID    string
		Account    string
		Method     string
		Route      string
		URI        string
		StatusCode int
		Consuming  int
		Type       int
		Connecting uint32
	}
)

// NewStats create a new stats middleware
func NewStats(config StatsConfig) echo.MiddlewareFunc {
	var connectingCount uint32
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			atomic.AddUint32(&connectingCount, 1)
			defer atomic.AddUint32(&connectingCount, ^uint32(0))
			startedAt := time.Now().UnixNano()
			req := c.Request()
			uri, _ := url.QueryUnescape(req.RequestURI)
			if uri == "" {
				uri = req.RequestURI
			}
			info := &StatsInfo{
				RequestID:  context.GetRequestID(c),
				Method:     req.Method,
				Route:      c.Path(),
				URI:        uri,
				Connecting: connectingCount,
				IP:         c.RealIP(),
				TrackID:    context.GetTrackID(c),
			}
			err = next(c)

			consuming := int(time.Now().UnixNano()-startedAt) / int(time.Millisecond)
			info.Consuming = consuming

			if err != nil {
				info.StatusCode = xerror.GetStatusCode(err)
			} else {
				info.StatusCode = context.GetStatus(c)
			}
			info.Type = info.StatusCode / 100
			us := context.GetUserSession(c)
			if us != nil {
				info.Account = us.GetAccount()
			}

			if config.OnStats != nil {
				config.OnStats(info)
			}
			return
		}
	}
}
