package controller

import (
	"go.uber.org/zap"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
)

type (
	bookCtrl struct {
	}
	bookBatchAddParams struct {
		Category string `valid:"in(xBiQuGe)"`
		Start    int    `valid:"xIntRange(1|100000)"`
		End      int    `valid:"xIntRange(1|100000)"`
	}
)

func init() {
	ctrl := bookCtrl{}
	books := router.NewGroup("/books")

	books.Add(
		"POST",
		"/v1/batch-add",
		middleware.Session,
		middleware.IsSu,
		ctrl.batchAdd,
	)
}

// batchAdd 批量添加
func (c *bookCtrl) batchAdd(ctx iris.Context) {
	params := &bookBatchAddParams{}
	err := validate(params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	go func() {
		start := params.Start
		end := params.End
		category := params.Category
		userLogger := getUserLogger(ctx)
		if start >= end {
			return
		}
		limit := 3
		wait := make(chan bool, limit)
		for i := start; i < end; i++ {
			wait <- true
			id := i
			go func() {
				service.AddBook(category, id)
				<-wait
			}()
		}
		userLogger.Info("batch add books done",
			zap.Int("start", start),
			zap.Int("end", end),
		)
	}()
	resNoContent(ctx)
}
