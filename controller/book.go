package controller

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
)

const ()

type (
	bookCtrl        struct{}
	addSourceParams struct {
		Category string `json:"category,omitempty" valid:"in(xBiQuGe|biQuGe)"`
		ID       int    `json:"id,omitempty" valid:"xIntRange(1|100000)"`
	}
	batchAddSourceParams struct {
		Category string `json:"category,omitempty" valid:"in(xBiQuGe|biQuGe)"`
		Start    int    `json:"start,omitempty" valid:"xIntRange(1|100000)"`
		Limit    int    `json:"limit,omitempty" valid:"xIntRange(1|10000)"`
	}
)

func init() {
	books := router.NewGroup("/books")
	ctrl := bookCtrl{}
	books.Add("POST", "/v1/sources", ctrl.addSource)
	books.Add("POST", "/v1/batch-sources", ctrl.batchAddSource)
}

// addSource add book source
func (c *bookCtrl) addSource(ctx iris.Context) {
	params := addSourceParams{}
	err := validate(&params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	err = service.AddNovel(params.Category, params.ID)
	if err != nil {
		resErr(ctx, err)
		return
	}
	resCreated(ctx, nil)
}

// batchAddSource bacth add book source
func (c *bookCtrl) batchAddSource(ctx iris.Context) {
	params := batchAddSourceParams{}
	err := validate(&params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	for index := 0; index < params.Limit; index++ {
		err = service.AddNovel(params.Category, params.Start+index)
		if err != nil {
			break
		}
	}
	if err != nil {
		resErr(ctx, err)
		return
	}
	resCreated(ctx, nil)
}
