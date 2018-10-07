package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/novel/xerror"

	"github.com/labstack/echo"
)

func TestNewStats(t *testing.T) {
	t.Run("handle success", func(t *testing.T) {
		done := false
		fn := NewStats(StatsConfig{
			OnStats: func(info *StatsInfo) {
				if info.StatusCode != http.StatusNoContent {
					t.Fatalf("stats info is wrong")
				}
				done = true
			},
		})(func(c echo.Context) error {
			return nil
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		if err != nil {
			t.Fatalf("stats midlleware fail, %v", err)
		}
		if !done {
			t.Fatalf("on stats is not called")
		}
	})

	t.Run("handle error", func(t *testing.T) {
		done := false
		fn := NewStats(StatsConfig{
			OnStats: func(info *StatsInfo) {
				if info.StatusCode != http.StatusInternalServerError {
					t.Fatalf("stats info is wrong")
				}
				done = true
			},
		})(func(c echo.Context) error {
			return errors.New("abcd")
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		if err == nil {
			t.Fatalf("stats should return error")
		}
		if !done {
			t.Fatalf("on stats is not called")
		}
	})

	t.Run("handle http error", func(t *testing.T) {
		done := false
		fn := NewStats(StatsConfig{
			OnStats: func(info *StatsInfo) {
				if info.StatusCode != http.StatusForbidden {
					t.Fatalf("stats info is wrong")
				}
				done = true
			},
		})(func(c echo.Context) error {
			return &xerror.HTTPError{
				Message:    "abc",
				StatusCode: http.StatusForbidden,
			}
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		if err == nil {
			t.Fatalf("stats should return error")
		}
		if !done {
			t.Fatalf("on stats is not called")
		}
	})
}
