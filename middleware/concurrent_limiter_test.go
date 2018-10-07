package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
)

func TestConcurrentLimiter(t *testing.T) {
	e := echo.New()
	t.Run("reset after done", func(t *testing.T) {
		fn := NewConcurrentLimiter(ConcurrentLimiterConfig{
			Category: "reset-true",
			Keys: []string{
				":ip",
				"h:X-Token",
				"q:id",
				"p:user",
				"account",
			},
			TTL:   time.Second,
			Reset: true,
		})(func(c echo.Context) error {
			return nil
		})
		buf := []byte(`{
			"account": "vicanso"
		}`)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/?id=1", nil)
		r.Header.Set("X-Token", "token")
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		c.SetParamNames("user")
		c.SetParamValues("tree.xie")
		context.SetRequestBody(c, buf)
		err := fn(c)
		if err != nil {
			t.Fatalf("concurrent limit(reset:true) fail, %v", err)
		}

		// 因为是reset:true，在上一次处理完成后，再次调用也可正常执行
		c = e.NewContext(r, w)
		err = fn(c)
		if err != nil {
			t.Fatalf("concurrent limit(reset:true) twice fail, %v", err)
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
		})(func(c echo.Context) error {
			return nil
		})
		buf := []byte(`{
			"account": "vicanso"
		}`)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/?id=1", nil)
		r.Header.Set("X-Token", "token")
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		context.SetRequestBody(c, buf)
		err := fn(c)
		if err != nil {
			t.Fatalf("concurrent limit(reset:false) fail, %v", err)
		}

		// 因为是reset:false，在上一次处理完成后，未过期时再次调用会出错
		err = fn(c)
		if err != errSubmitTooFrequently {
			t.Fatalf("concurrent limit(reset:false) should be fail")
		}

		time.Sleep(1100 * time.Millisecond)
		// 在过期时间后，则可正常执行
		err = fn(c)
		if err != nil {
			t.Fatalf("concurrent limit(reset:false) expired should be success")
		}
	})
}
