package controller

import (
	"strconv"

	"go.uber.org/zap"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
)

type (
	// BookCtrl book controller
	BookCtrl struct {
	}
	// BookBatchAddParams params for book batch add
	BookBatchAddParams struct {
		Category string `valid:"in(xBiQuGe)"`
		Start    int    `valid:"xIntRange(1|100000)"`
		End      int    `valid:"xIntRange(1|100000)"`
	}
	// BookUpdateChaptersParams params for book update chapters
	BookUpdateChaptersParams struct {
		Author string `valid:"runelength(1|64)"`
		Name   string `valid:"runelength(1|64)"`
		Limit  int    `valid:"xIntRange(1|10)"`
	}
	// BookQueryOptionsParams params for the query
	BookQueryOptionsParams struct {
		Limit  string `json:"limit,omitempty" valid:"range(1|20)"`
		Offset string `json:"offset,omitempty" valid:"numeric"`
		Field  string `json:"field,omitempty" valid:"runelength(2|64)"`
		Order  string `json:"order,omitempty" valid:"runelength(2|32)"`
	}
)

func init() {
	ctrl := BookCtrl{}
	books := router.NewGroup("/books")
	// 书籍查询
	books.Add(
		"GET",
		"/v1",
		ctrl.list,
	)
	// 章节查询
	books.Add(
		"GET",
		"/v1/:book/chapters",
		ctrl.listChapter,
	)

	// 批量新增
	books.Add(
		"POST",
		"/v1/batch-add",
		newDefaultTracker(cs.ActionBookBatchAdd, nil),
		middleware.Session,
		middleware.IsSu,
		ctrl.batchAdd,
	)
	// 更新章节
	books.Add(
		"PATCH",
		"/v1/chapters",
		newDefaultTracker(cs.ActionBookUpdateChapters, nil),
		middleware.Session,
		middleware.IsSu,
		ctrl.updateChapters,
	)
}

// batchAdd 批量添加
func (c *BookCtrl) batchAdd(ctx iris.Context) {
	params := &BookBatchAddParams{}
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

// updateChapters update chapters
func (c *BookCtrl) updateChapters(ctx iris.Context) {
	params := &BookUpdateChaptersParams{}
	err := validate(params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	service.UpdateBookChapter(params.Author, params.Name, params.Limit)
}

func getQueryOptions(ctx iris.Context) (opts *model.QueryOptions, err error) {
	params := &BookQueryOptionsParams{}
	err = validate(params, getRequestQuery(ctx))
	if err != nil {
		return
	}
	limit, _ := strconv.Atoi(params.Limit)
	offset, _ := strconv.Atoi(params.Offset)
	opts = &model.QueryOptions{
		Limit:  limit,
		Offset: offset,
		Order:  params.Order,
		Field:  params.Field,
	}
	return
}

// list list the book
func (c *BookCtrl) list(ctx iris.Context) {
	opts, err := getQueryOptions(ctx)
	if err != nil {
		resErr(ctx, err)
		return
	}

	query := getRequestQuery(ctx)
	q := query["q"]
	if q != "" {
		books, err := service.ListBookByKeyword(q, opts)
		if err != nil {
			resErr(ctx, err)
			return
		}
		setCache(ctx, "1m")
		res(ctx, map[string]interface{}{
			"books": books,
		})
		return
	}

	conditions := model.Book{}
	count, err := service.CountBook(conditions)
	if err != nil {
		resErr(ctx, err)
		return
	}
	books, err := service.ListBook(&conditions, opts)
	if err != nil {
		resErr(ctx, err)
		return
	}
	setCache(ctx, "5m")
	res(ctx, map[string]interface{}{
		"books": books,
		"count": count,
	})
}

// listChapter list the chapter
func (c *BookCtrl) listChapter(ctx iris.Context) {
	opts, err := getQueryOptions(ctx)
	if err != nil {
		resErr(ctx, err)
		return
	}
	bookID, err := ctx.Params().GetInt("book")
	if err != nil {
		resErr(ctx, err)
		return
	}
	conditions := &model.Chapter{
		BookID: uint(bookID),
	}
	chapters, err := service.ListBookChapters(conditions, opts)
	if err != nil {
		resErr(ctx, err)
		return
	}
	setCache(ctx, "5m")
	res(ctx, map[string]interface{}{
		"chapters": chapters,
	})
}
