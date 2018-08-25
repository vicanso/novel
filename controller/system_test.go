package controller

import (
	"sort"
	"strings"
	"testing"

	"github.com/kataras/iris"

	"github.com/vicanso/novel/utils"
)

func TestSystemCtrl(t *testing.T) {
	ctrl := systemCtrl{}

	t.Run("getStatus", func(t *testing.T) {
		ctx := utils.NewResContext()
		ctrl.getStatus(ctx)
		data := utils.GetBody(ctx).(iris.Map)
		keys := []string{}
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if strings.Join(keys, ",") != "goMaxProcs,startedAt,status,uptime,version" {
			t.Fatalf("get status fail")
		}
	})

	t.Run("getStats", func(t *testing.T) {
		ctx := utils.NewResContext()
		ctrl.getStats(ctx)
		data := utils.GetBody(ctx).(iris.Map)
		keys := []string{}
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if strings.Join(keys, ",") != "heapInuse,heapSys,routineCount,sys" {
			t.Fatalf("get stats fail")
		}
	})
}
