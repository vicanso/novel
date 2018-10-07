package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/vicanso/novel/xerror"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/service"
)

var (
	// submit too frequently
	errSubmitTooFrequently = &xerror.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   xerror.ErrCategoryValidte,
		Message:    "submit too frequently",
	}
)

type (
	// ConcurrentLimiterConfig concurrent limiter config
	ConcurrentLimiterConfig struct {
		Category string
		Keys     []string
		TTL      time.Duration
		// 是否在完成后将限制重置
		Reset bool
	}
	// ConcurrentKeyInfo the concurrent key's info
	ConcurrentKeyInfo struct {
		Name   string
		Params bool
		Query  bool
		Header bool
		Body   bool
		IP     bool
	}
)

// NewConcurrentLimiter create a concurrent limitter middleware
func NewConcurrentLimiter(config ConcurrentLimiterConfig) echo.MiddlewareFunc {
	keys := make([]*ConcurrentKeyInfo, 0)
	// 根据配置生成key的处理
	for _, key := range config.Keys {
		if key == ":ip" {
			keys = append(keys, &ConcurrentKeyInfo{
				IP: true,
			})
			continue
		}
		if strings.HasPrefix(key, "h:") {
			keys = append(keys, &ConcurrentKeyInfo{
				Name:   key[2:],
				Header: true,
			})
			continue
		}
		if strings.HasPrefix(key, "q:") {
			keys = append(keys, &ConcurrentKeyInfo{
				Name:  key[2:],
				Query: true,
			})
			continue
		}
		if strings.HasPrefix(key, "p:") {
			keys = append(keys, &ConcurrentKeyInfo{
				Name:   key[2:],
				Params: true,
			})
			continue
		}
		keys = append(keys, &ConcurrentKeyInfo{
			Name: key,
			Body: true,
		})
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			values := make([]string, len(keys)+1)
			values[0] = config.Category
			req := c.Request()
			// 获取lock的key
			for i, key := range keys {
				v := ""
				name := key.Name
				if key.IP {
					v = c.RealIP()
				} else if key.Header {
					v = req.Header.Get(name)
				} else if key.Query {
					query := context.GetRequestQuery(c)
					v = query[name]
				} else if key.Params {
					v = c.Param(name)
				} else {
					body := context.GetRequestBody(c)
					v = json.Get(body, name).ToString()
				}
				values[i+1] = v
			}
			lockKey := strings.Join(values, ",")
			// 判断是否可以lock成功
			success, done, err := service.LockWithDone(lockKey, config.TTL)
			if err != nil {
				return
			}
			// 如果lock失败，则出错
			if !success {
				err = errSubmitTooFrequently
				return
			}
			// 如果设置了完成后重置锁，则重置
			if config.Reset {
				defer done()
			}
			return next(c)
		}
	}
}
