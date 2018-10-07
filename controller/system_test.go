package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/vicanso/novel/context"
)

func TestSystemCtrl(t *testing.T) {
	ctrl := SystemCtrl{}
	e := echo.New()

	t.Run("getStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := e.NewContext(nil, w)
		err := ctrl.getStatus(c)
		if err != nil {
			t.Fatalf("get status fail, %v", err)
		}
		data := context.GetBody(c).(*StatusRes)
		if data.Pid == 0 {
			t.Fatalf("get status fail")
		}
	})

	t.Run("getStats", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := e.NewContext(nil, w)
		err := ctrl.getStats(c)
		if err != nil {
			t.Fatalf("get stats fail, %v", err)
		}
		data := context.GetBody(c).(*StatsRes)
		if data.Sys == 0 {
			t.Fatalf("get stats fail")
		}
	})

	t.Run("get routes", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := e.NewContext(nil, w)
		err := ctrl.getRoutes(c)
		if err != nil {
			t.Fatalf("get routes fail, %v", err)
		}
		data := context.GetBody(c).(*RoutesRes)
		if data == nil {
			t.Fatalf("get routes fail")
		}
	})

	t.Run("get route count", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := e.NewContext(nil, w)
		err := ctrl.getRouteCounts(c)
		if err != nil {
			t.Fatalf("get route count fail, %v", err)
		}
		data := context.GetBody(c).(map[string]interface{})
		if data["createdAt"] == nil || data["counts"] == nil {
			t.Fatalf("get route count fail")
		}
	})
}
