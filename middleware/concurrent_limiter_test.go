package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vicanso/novel/util"
)

func TestConcurrentLimiter(t *testing.T) {
	t.Run("reset after done", func(t *testing.T) {
		fn := NewConcurrentLimiter(ConcurrentLimiterConfig{
			Category: "reset-true",
			Keys: []string{
				"h:X-Token",
				"q:id",
				"account",
			},
			TTL:   time.Second,
			Reset: true,
		})
		buf := []byte(`{
			"account": "vicanso"
		}`)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/?id=1", nil)
		r.Header.Set("X-Token", "token")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("concurrent limit(reset:true) fail")
		}

		// 因为是reset:true，在上一次处理完成后，再次调用也可正常执行
		ctx = util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("concurrent limit(reset:true) twice fail")
		}
	})

	t.Run("not reset after done", func(t *testing.T) {
		fn := NewConcurrentLimiter(ConcurrentLimiterConfig{
			Category: "reset-false",
			Keys: []string{
				"h:X-Token",
				"q:id",
				"account",
			},
			TTL:   time.Second,
			Reset: false,
		})
		buf := []byte(`{
			"account": "vicanso"
		}`)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/?id=1", nil)
		r.Header.Set("X-Token", "token")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("concurrent limit(reset:false) fail")
		}

		// 因为是reset:false，在上一次处理完成后，未过期时再次调用会出错
		ctx = util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		fn(ctx)
		if ctx.GetStatusCode() == http.StatusOK {
			t.Fatalf("concurrent limit(reset:false) should be fail")
		}
		time.Sleep(1100 * time.Millisecond)

		// 在过期时间后，则可正常执行
		ctx = util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("concurrent limit(reset:false) expired should be success")
		}
	})
}
