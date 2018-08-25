package controller

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
)

type (
	novelCtrl       struct{}
	addSourceParams struct {
		Category string `json:"category,omitempty" valid:"in(xBiQuGe|biQuGe)"`
		ID       int    `json:"id,omitempty" valid:"-"`
	}
)

func init() {
	novels := router.NewGroup("/novels")
	ctrl := novelCtrl{}
	novels.Add("POST", "/v1/sources", ctrl.addSource)
}

// addSource add novel source
func (c *novelCtrl) addSource(ctx iris.Context) {
	params := &addSourceParams{}
	err := validate(params, getRequestBody(ctx))
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
