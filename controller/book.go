package controller

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/validate"
	"go.uber.org/zap"
)

const (
	bookAddFavCategory = "add"
)

type (
	// BookCtrl book controller
	BookCtrl struct{}
	// BookBatchAddParams params for book batch add
	BookBatchAddParams struct {
		Category string `valid:"in(xBiQuGe)"`
		Start    int    `valid:"xIntRange(1|100000)"`
		End      int    `valid:"xIntRange(1|100000)"`
	}
	// BookUpdateChaptersParams params for book update chapters
	BookUpdateChaptersParams struct {
		Limit int `valid:"xIntRange(1|10)"`
	}
	// UserActionParams user action params
	UserActionParams struct {
		Type string `valid:"in(like|view)"`
	}
	// BookFavToggleParams params for book fav toggle
	BookFavToggleParams struct {
		Category string `valid:"in(add|remove)"`
	}
)

func init() {
	books := router.NewGroup("/books")
	ctrl := BookCtrl{}
	isSu := middleware.IsSu

	// list the books
	books.Add(
		"GET",
		"/v1",
		ctrl.list,
	)

	// get the book's categories
	books.Add(
		"GET",
		"/v1/categories",
		ctrl.getCategories,
	)

	// update the book's info
	books.Add(
		"PATCH",
		"/v1/:id",
		ctrl.updateInfo,
		createTracker(cs.ActionBookUpdateInfo),
		userSession,
		isSu,
	)

	// batch add books
	books.Add(
		"POST",
		"/v1/batch-add",
		ctrl.batchAdd,
		createTracker(cs.ActionBookBatchAdd),
		userSession,
		isSu,
	)

	// update the book's chapters
	books.Add(
		"PATCH",
		"/v1/chapters/:id",
		ctrl.updateChapters,
		createTracker(cs.ActionBookUpdateChapters),
		userSession,
		isSu,
		middleware.NewConcurrentLimiter(middleware.ConcurrentLimiterConfig{
			Category: cs.ActionBookUpdateChapters,
			Keys: []string{
				"p:id",
			},
			// 只允许5分钟执行一次主动更新
			TTL: 300 * time.Second,
		}),
	)

	// list the book's chapters
	books.Add(
		"GET",
		"/v1/chapters/:id",
		ctrl.listChapaters,
	)

	// update the book's cover
	books.Add(
		"PATCH",
		"/v1/covers/:id",
		ctrl.updateCover,
		userSession,
		isSu,
	)

	// userAction the book
	books.Add(
		"POST",
		"/v1/actions/:id",
		ctrl.userAction,
		createTracker(cs.ActionUserBookAction),
		userSession,
		middleware.IsLogined,
	)

	books.Add(
		"GET",
		"/v1/recommends/:id",
		ctrl.getRecommendByID,
	)

	// 用户收藏
	books.Add(
		"GET",
		"/v1/favorites",
		ctrl.listFav,
		userSession,
		middleware.IsLogined,
	)
	// 用户收藏切换（添加删除）
	books.Add(
		"POST",
		"/v1/favorites/:id",
		ctrl.toggleFav,
		createTracker(cs.ActionUserFavToggle),
		userSession,
		middleware.IsLogined,
	)

	// get the book's info
	books.Add(
		"GET",
		"/v1/:id",
		ctrl.getInfo,
	)

}

// batchAdd batch add books
func (bc *BookCtrl) batchAdd(c echo.Context) (err error) {
	params := &BookBatchAddParams{}
	err = validate.Do(params, getRequestBody(c))
	if err != nil {
		return
	}
	logger := getContextLogger(c)
	go func() {
		start := params.Start
		end := params.End
		category := params.Category
		if start >= end {
			return
		}
		limit := 3
		wait := make(chan bool, limit)
		for i := start; i < end; i++ {
			wait <- true
			id := i
			go func() {
				bookService.Add(category, id)
				<-wait
			}()
		}
		logger.Info("batch add books done",
			zap.Int("start", start),
			zap.Int("end", end),
		)
	}()
	return
}

// updateChapters update the latest book's chapters
func (bc *BookCtrl) updateChapters(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := &BookUpdateChaptersParams{}
	err = validate.Do(params, getRequestBody(c))
	if err != nil {
		return
	}
	err = bookService.UpdateChapters(id)
	return
}

// updateCover update book's cover
func (bc *BookCtrl) updateCover(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	cover, err := bookService.UpdateCover(id)
	if err != nil {
		return
	}
	res(c, map[string]string{
		"cover": cover,
	})
	return
}

// list list the book
func (bc *BookCtrl) list(c echo.Context) (err error) {
	params := &service.BookQueryParams{}
	err = validate.Do(params, getRequestQuery(c))
	if err != nil {
		return
	}
	books, err := bookService.List(params)
	if err != nil {
		return
	}
	m := map[string]interface{}{
		"books": books,
	}
	offset := params.Offset
	if offset == "0" || offset == "" {
		count, err := bookService.Count(params)
		if err != nil {
			return err
		}
		m["count"] = count
	}
	setCache(c, "10s")
	res(c, m)
	return
}

// get get the book's info
func (bc *BookCtrl) getInfo(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	book, err := bookService.GetInfo(id)
	if err != nil {
		return
	}
	latestChapter, err := bookService.GetLatestChapter(id, "title,updatedAt")
	if err != nil {
		return
	}
	chapterCount, err := bookService.CountChapters(id)
	if err != nil {
		return
	}
	data := map[string]interface{}{
		"book":          book,
		"chapterCount":  chapterCount,
		"latestChapter": latestChapter,
	}
	setCacheWithSMaxAge(c, "5m", "30s")
	res(c, data)
	return
}

// update update the book's info
func (bc *BookCtrl) updateInfo(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := &service.BookUpdateParams{}
	err = validate.Do(params, getRequestBody(c))
	if err != nil {
		return
	}
	err = bookService.UpdateInfo(id, params)
	if err == nil && params.SourceCover != "" {
		bookService.UpdateCover(id)
	}
	return
}

// listChapaters list the book's chapters
func (bc *BookCtrl) listChapaters(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := &service.BookChapterQueryParams{}
	err = validate.Do(params, getRequestQuery(c))
	if err != nil {
		return
	}
	chapters, err := bookService.ListChapters(id, params)
	if err != nil {
		return
	}
	m := make(map[string]interface{})
	m["chapters"] = chapters
	offset := params.Offset
	if offset == "0" || offset == "" {
		count, err := bookService.CountChapters(id)
		if err != nil {
			return err
		}
		m["count"] = count
	}
	setCache(c, "5m")
	res(c, m)
	return
}

// userAction user actions
func (bc *BookCtrl) userAction(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	parmas := &UserActionParams{}
	err = validate.Do(parmas, getRequestBody(c))
	if err != nil {
		return
	}
	if parmas.Type == "view" {
		err = bookService.View(id)
	} else {
		err = bookService.Like(id)
	}
	return
}

// getCategories get the book's categories
func (bc *BookCtrl) getCategories(c echo.Context) (err error) {
	categories := make(map[string]int)
	_, err = service.RedisGet(cs.CacheBookCategories, &categories)
	if err != nil && !service.IsRedisNil(err) {
		return
	}
	// 因为有可能为redis nil error，所以重置
	err = nil
	setCache(c, "5m")
	res(c, map[string]interface{}{
		"categories": categories,
	})
	return
}

// getRecommendByID get the recommend by book's id
func (bc *BookCtrl) getRecommendByID(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	parmas := &service.BookQueryParams{}
	err = validate.Do(parmas, getRequestQuery(c))
	if err != nil {
		return
	}
	info, err := bookService.GetInfo(id)
	if err != nil {
		return
	}
	rand.Seed(time.Now().UnixNano())
	category := info.Category[rand.Intn(len(info.Category))]

	parmas.Author = info.Author
	authorResult, err := bookService.List(parmas)
	if err != nil {
		return
	}

	parmas.Author = ""
	parmas.Category = category
	catResult, err := bookService.List(parmas)
	if err != nil {
		return
	}
	result := make([]*model.Book, 0)
	for _, item := range authorResult {
		if item.ID != uint(id) {
			result = append(result, item)
		}
	}
	for _, item := range catResult {
		if item.ID != uint(id) {
			// 有可能与通过作者查询的重复
			found := false
			for _, tmp := range result {
				if tmp.ID == item.ID {
					found = true
				}
			}
			if !found {
				result = append(result, item)
			}
		}
	}
	setCache(c, "1m")
	res(c, map[string]interface{}{
		"books": result,
	})
	return
}

// toggleFav toggle the user fav
func (bc *BookCtrl) toggleFav(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := &BookFavToggleParams{}
	err = validate.Do(params, getRequestBody(c))
	if err != nil {
		return
	}
	account := context.GetUserSession(c).GetAccount()
	if params.Category == bookAddFavCategory {
		err = bookService.AddFav(account, id)
	} else {
		err = bookService.RemoveFav(account, id)
	}
	if err != nil {
		return
	}
	return
}

// listFav list the user fav
func (bc *BookCtrl) listFav(c echo.Context) (err error) {
	account := context.GetUserSession(c).GetAccount()
	result, err := bookService.ListFav(account)
	if err != nil {
		return
	}
	res(c, map[string]interface{}{
		"favs": result,
	})
	return
}
