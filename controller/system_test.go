package controller

import (
	"net/http"
	"testing"

	"github.com/vicanso/novel/util"
)

func TestSystemCtrl(t *testing.T) {
	ctrl := systemCtrl{}

	t.Run("getStatus", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getStatus(ctx)
		data := util.GetBody(ctx).(*statusRes)
		if data.Pid == 0 {
			t.Fatalf("get status fail")
		}
	})

	t.Run("getStats", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getStats(ctx)
		data := util.GetBody(ctx).(*statsRes)
		if data.Sys == 0 {
			t.Fatalf("get stats fail")
		}
	})

	t.Run("get routes", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getRoutes(ctx)
		_ = util.GetBody(ctx).(*routesRes)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("get routes fail")
		}
	})

	t.Run("get route count", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getRouteCounts(ctx)
		data := util.GetBody(ctx).(map[string]interface{})
		if data["createdAt"] == nil || data["counts"] == nil {
			t.Fatalf("get route count fail")
		}
	})
}
