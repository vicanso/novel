package controller

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/model"
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
	updateBookParams struct {
		Status int    `json:"status,omitempty" valid:"xIntRange(1|3),optional"`
		Brief  string `json:"brief,omitempty" valid:"runelength(1|300),optional"`
	}
)

func init() {
	books := router.NewGroup("/books")
	ctrl := bookCtrl{}
	books.Add("POST", "/v1/sources", ctrl.addSource)
	books.Add("POST", "/v1/batch-sources", ctrl.batchAddSource)
	books.Add("GET", "/v1/:id", ctrl.detail)
	books.Add("PATCH", "/v1/:id", ctrl.update)
	books.Add("GET", "", ctrl.list)
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

// list get books
func (c *bookCtrl) list(ctx iris.Context) {
	params := &service.NovelQueryParams{}
	err := validate(params, getRequestQuery(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	books, err := service.ListNovel(params)
	if err != nil {
		resErr(ctx, err)
		return
	}
	data := iris.Map{
		"books": books,
	}
	if params.Offset == "" || params.Offset == "0" {
		count, err := service.CountNovel(params)
		if err != nil {
			resErr(ctx, err)
			return
		}
		data["count"] = count
	}
	setCache(ctx, "1m")
	res(ctx, data)
}

// detail get the datail info
func (c *bookCtrl) detail(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		resErr(ctx, err)
		return
	}
	b := &model.Book{}
	b.ID = uint(id)
	err = model.GetClient().Model(b).First(b).Error
	if err != nil {
		resErr(ctx, err)
		return
	}
	res(ctx, b)
}

// update update info
func (c *bookCtrl) update(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		resErr(ctx, err)
		return
	}
	params := &updateBookParams{}
	err = validate(params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	err = model.UpdateBookByID(uint(id), &model.Book{
		Status: params.Status,
		Brief:  params.Brief,
	})
	if err != nil {
		resErr(ctx, err)
		return
	}
	resNoContent(ctx)
}
