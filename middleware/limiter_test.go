package middleware

import (
	"testing"
	"time"

	"github.com/vicanso/novel/xerror"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/global"
)

func TestNewLimiter(t *testing.T) {
	t.Run("limit pass", func(t *testing.T) {
		fn := NewLimiter(LimiterConfig{
			Max: 1,
		})(func(c echo.Context) error {
			return nil
		})
		e := echo.New()
		c := e.NewContext(nil, nil)
		err := fn(c)
		if err != nil {
			t.Fatalf("pass limiter fail, %v", err)
		}
	})

	t.Run("over limit", func(t *testing.T) {
		global.StartApplication()
		fn := NewLimiter(LimiterConfig{
			Max: 0,
		})(func(c echo.Context) error {
			return nil
		})
		e := echo.New()
		c := e.NewContext(nil, nil)
		err := fn(c)
		if err == nil || err.(*xerror.HTTPError) != errTooManyRequest {
			t.Fatalf("should return error over limit")
		}

		if global.IsApplicationRunning() {
			t.Fatalf("the application should be paused after over limit")
		}
	})
}

func TestResetApplication(t *testing.T) {
	global.PauseApplication()
	if global.IsApplicationRunning() {
		t.Fatalf("application should be pause")
	}
	resetApplicationStatus(time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if !global.IsApplicationRunning() {
		t.Fatalf("application should resume to running")
	}
}
