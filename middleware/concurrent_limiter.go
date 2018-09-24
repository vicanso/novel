package middleware

import (
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
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
		Query  bool
		Header bool
		Body   bool
	}
)

// NewConcurrentLimiter 并发请求限制中间件
func NewConcurrentLimiter(conf ConcurrentLimiterConfig) iris.Handler {
	keys := make([]*ConcurrentKeyInfo, 0)
	for _, key := range conf.Keys {
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
		keys = append(keys, &ConcurrentKeyInfo{
			Name: key,
			Body: true,
		})
	}
	return func(ctx iris.Context) {
		values := make([]string, len(keys)+1)
		values[0] = conf.Category
		for i, key := range keys {
			v := ""
			name := key.Name
			if key.Header {
				v = ctx.GetHeader(name)
			} else if key.Query {
				query := util.GetRequestQuery(ctx)
				v = query[name]
			} else {
				body := util.GetRequestBody(ctx)
				v = json.Get(body, name).ToString()
			}
			values[i+1] = v
		}
		lockKey := strings.Join(values, ",")
		success, done, err := service.LockWithDone(lockKey, conf.TTL)
		if err != nil {
			resErr(ctx, err)
			return
		}
		if !success {
			resErr(ctx, util.ErrSubmitTooFrequently)
			return
		}
		if conf.Reset {
			defer done()
		}
		ctx.Next()
	}
}
