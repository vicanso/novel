package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
)

func TestNewTracker(t *testing.T) {
	done := false
	onTrack := func(info *TrackerInfo, c echo.Context) {
		if info.Result != HandleSuccess {
			t.Fatalf("tracker info fail")
		}
		done = true
	}
	fn := NewTracker(TrackerConfig{
		OnTrack: onTrack,
	})(func(c echo.Context) error {
		return nil
	})
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/?c=3", nil)
	w := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetParamNames("a")
	c.SetParamValues("1")
	context.SetRequestBody(c, []byte(`{"b": "2"}`))
	err := fn(c)
	if err != nil {
		t.Fatalf("tracker middleware fail, %v", err)
	}
	if !done {
		t.Fatalf("on track is not called")
	}
}
