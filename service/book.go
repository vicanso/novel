package service

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/vicanso/novel-spider/mq"
	"github.com/vicanso/novel-spider/novel"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/helper"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/xlog"
	"go.uber.org/zap"
)

var (
	mqClient    *mq.MQ
	uploadURL   string
	uploadToken string
)

const (
	bookCoverCategory = "book-cover"
	bookViewCount     = "view_count"
	bookLikeCount     = "like_count"
	bookWordCount     = "word_count"
	bookCategory      = "category"
	bookAuthor        = "author"
	bookName          = "name"
	bookStatus        = "status"
)

type (
	// Book book service struct
	Book struct {
	}
	// ChapterCountResult chapter count's result
	ChapterCountResult struct {
		Total int
	}

	// BookQueryParams params for the query
	BookQueryParams struct {
		Limit    string `json:"limit,omitempty" valid:"range(1|20)"`
		Offset   string `json:"offset,omitempty" valid:"numeric"`
		Field    string `json:"field,omitempty" valid:"runelength(2|64)"`
		Order    string `json:"order,omitempty" valid:"runelength(2|32)"`
		Q        string `json:"q,omitempty" valid:"runelength(1|32),optional"`
		Category string `json:"category,omitempty" valid:"runelength(2|8),optional"`
		Status   string `json:"status,omitempty" valid:"in(0|1|2),optional"`
	}
	// BookUpdateParams params for update
	BookUpdateParams struct {
		Brief    string `json:"brief,omitempty" valid:"runelength(5|2000),optional"`
		Status   int    `json:"status,omitempty" valid:"xIntIn(1|2),optional"`
		Category string `json:"category,omitempty" valid:"runelength(2|32),optional"`
	}
	// BookChapterQueryParams params for the query
	BookChapterQueryParams struct {
		Limit  string `json:"limit,omitempty" valid:"range(1|20)"`
		Offset string `json:"offset,omitempty" valid:"numeric"`
		Field  string `json:"field,omitempty" valid:"runelength(2|64)"`
		Order  string `json:"order,omitempty" valid:"runelength(2|32)"`
	}
)

func init() {
	address := config.GetStringSlice("nsq.lookup.address")

	c := &mq.MQ{
		LookupAddress: address,
		Logger:        xlog.Logger(),
	}
	c.FreshNodes()
	go c.TimedFreshNodes(time.Second * 60)

	if err := initReceiveBasicInfoEvent(c); err != nil {
		panic(err)
	}
	if err := initReceiveChapterEvent(c); err != nil {
		panic(err)
	}

	uploadURL = config.GetString("tiny.host") + config.GetString("tiny.upload")
	uploadToken = config.GetString("tiny.token")

	mqClient = c
}

func initReceiveBasicInfoEvent(c *mq.MQ) (err error) {
	logger := xlog.Logger()
	cb := func(info *novel.BasicInfo) {
		if info == nil {
			return
		}
		b := &model.Book{
			Name:   info.Name,
			Author: info.Author,
			Brief:  info.Brief,
			Category: []string{
				info.Category,
			},
			Source:      info.Source,
			SourceID:    info.SourceID,
			SourceCover: info.Cover,
		}
		err = getClient().Create(b).Error
		if err != nil {
			logger.Error("add book fail",
				zap.String("name", info.Name),
				zap.String("author", info.Author),
				zap.Error(err),
			)
		}
	}
	err = c.SubReceiveNovel(cb)
	return
}

func initReceiveChapterEvent(c *mq.MQ) (err error) {
	logger := xlog.Logger()
	cb := func(c *novel.Chapter) {
		// 小于1000字的章节认为非正常数据
		wordCount := len(c.Content)
		if wordCount < 1000 {
			return
		}
		b := &model.Book{}
		err := getClient().
			Where(&model.Book{
				Name:   c.Name,
				Author: c.Author,
			}).
			First(b).Error
		if err != nil || b.ID == 0 {
			return
		}
		count := 0
		getClient().
			Model(&model.Chapter{}).
			Where(&model.Chapter{
				Index:  c.Index,
				BookID: b.ID,
			}).
			Count(&count)
		// 如果章节已存在，则跳过
		if count != 0 {
			return
		}
		chapter := &model.Chapter{
			Title:     c.Title,
			Content:   c.Content,
			WordCount: wordCount,
			Index:     c.Index,
			BookID:    b.ID,
		}
		err = getClient().Create(chapter).Error
		if err != nil {
			logger.Error("update chpater fail",
				zap.String("author", c.Author),
				zap.String("name", c.Name),
				zap.Int("index", c.Index),
			)
		}
		bookService := &Book{}
		bookService.UpdateWordCount(int(b.ID))
	}
	err = c.SubReceiveChapter(cb)
	return
}

// Add add the book
func (b *Book) Add(category string, id int) (err error) {
	count := 0
	err = getClient().
		Model(&model.Book{}).
		Where(&model.Book{
			Source:   category,
			SourceID: id,
		}).
		Count(&count).Error
	if err != nil || count != 0 {
		return
	}

	return mqClient.Pub(mq.TopicAddNovel, &novel.Source{
		Category: category,
		ID:       id,
	})
}

// UpdateChapters update book's chapters
func (b *Book) UpdateChapters(id, limit int) (err error) {
	book := &model.Book{}
	book.ID = uint(id)
	err = getClient().
		Where(book).
		First(book).
		Error
	if err != nil {
		return
	}
	chapters := make([]*model.Chapter, 0)
	err = getClientByOptions(&model.QueryOptions{
		Field: "index,id",
		Order: "-id",
		Limit: limit,
	}, nil).
		Where(&model.Chapter{
			BookID: book.ID,
		}).
		Find(&chapters).
		Error
	if err != nil {
		return
	}

	start := 0
	if len(chapters) != 0 {
		for i, j := 0, len(chapters)-1; i < j; i, j = i+1, j-1 {
			chapters[i], chapters[j] = chapters[j], chapters[i]
		}
		for _, chapter := range chapters {
			if start == 0 {
				start = chapter.Index
				continue
			}
			// 如果下一章也存在，最新章节修改为下一章
			if chapter.Index-start == 1 {
				start = chapter.Index
			}
		}
		// 需要更新的章节为下一章
		start++
	}

	err = b.AddChapter(book.Source, book.SourceID, start)
	return
}

// AddChapter add chapter
func (b *Book) AddChapter(category string, id, latestChapter int) (err error) {
	return mqClient.Pub(mq.TopicUpdateChapter, &novel.Source{
		Category:      category,
		ID:            id,
		LatestChapter: latestChapter,
	})
}

// UpdateCover update cover
func (b *Book) UpdateCover(id int) (err error) {
	book := &model.Book{}
	book.ID = uint(id)
	err = getClient().Where(book).First(book).Error
	buf, err := helper.HTTPGet(book.SourceCover, nil)
	if err != nil {
		return
	}

	data := make(map[string]interface{})
	data["token"] = uploadToken
	data["category"] = bookCoverCategory
	data["fileType"] = "jpeg"
	data["maxAge"] = "720h"
	data["data"] = base64.StdEncoding.EncodeToString(buf)
	res, err := helper.HTTPPost(uploadURL, data, nil)
	if err != nil {
		return
	}
	file := json.Get(res, "file").ToString()
	err = getClient().Model(book).Update(&model.Book{
		Cover: file,
	}).Error
	return
}

// bookKeywordSearch the keyword search
func bookKeywordSearch(client *gorm.DB, q string) *gorm.DB {
	if q != "" {
		key := "%" + q + "%"
		client = client.Where(bookName+" LIKE ?", key).
			Or(bookAuthor+" LIKE ?", key)
	}
	return client
}

// getWhereConditions get where conditions
func getWhereConditions(params *BookQueryParams) (query string, args []interface{}) {
	args = make([]interface{}, 0)
	sql := []string{}
	q := params.Q
	if q != "" {
		key := "%" + q + "%"
		str := fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", bookName, bookAuthor)
		sql = append(sql, str)
		args = append(args, key)
		args = append(args, key)
	}
	category := params.Category
	if category != "" {
		sql = append(sql, "? = ANY("+bookCategory+")")
		args = append(args, category)
	}
	status := params.Status
	if status != "" {
		sql = append(sql, bookStatus+" = ?")
		args = append(args, status)
	}
	query = strings.Join(sql, " AND ")
	return
}

// bookCategorySearch the category search
func bookCategorySearch(client *gorm.DB, category string) *gorm.DB {
	if category != "" {
		client = client.Where("? = ANY("+bookCategory+")", category)
	}
	return client
}

// List list the book
func (b *Book) List(params *BookQueryParams) (books []*model.Book, err error) {
	limit, _ := strconv.Atoi(params.Limit)
	offset, _ := strconv.Atoi(params.Offset)
	options := &model.QueryOptions{
		Limit:  limit,
		Offset: offset,
		Field:  params.Field,
		Order:  params.Order,
	}
	books = make([]*model.Book, 0)
	client := getClientByOptions(options, nil)

	query, args := getWhereConditions(params)
	client = client.Where(query, args...)

	err = client.Find(&books).Error
	return
}

// Count count the book
func (b *Book) Count(params *BookQueryParams) (count int, err error) {
	client := getClient().Model(&model.Book{})

	query, args := getWhereConditions(params)
	client = client.Where(query, args...)

	err = client.Count(&count).Error
	return
}

// GetInfo get the book's info
func (b *Book) GetInfo(id int) (book *model.Book, err error) {
	book = &model.Book{}
	book.ID = uint(id)
	err = getClient().Where(book).First(book).Error
	return
}

// UpdateInfo update the book's info
func (b *Book) UpdateInfo(id int, params *BookUpdateParams) (err error) {
	book := &model.Book{}
	book.ID = uint(id)
	updateInfo := &model.Book{}
	if params.Brief != "" {
		updateInfo.Brief = params.Brief
	}
	if params.Category != "" {
		updateInfo.Category = strings.Split(params.Category, ",")
	}
	if params.Status != 0 {
		updateInfo.Status = params.Status
	}
	err = getClient().Model(book).Updates(updateInfo).Error
	return
}

// UpdateWordCount update the book's word count
func (b *Book) UpdateWordCount(id int) (err error) {
	result := &ChapterCountResult{}
	err = getClient().
		Model(&model.Chapter{}).
		Select("sum(" + bookWordCount + ") as total").
		Where(&model.Chapter{
			BookID: uint(id),
		}).
		Scan(result).Error
	if err != nil || result.Total == 0 {
		return
	}
	book := &model.Book{}
	book.ID = uint(id)
	err = getClient().
		Where(book).
		First(book).Error
	if err != nil || book.WordCount == result.Total {
		return
	}

	err = getClient().Model(book).Updates(&model.Book{
		WordCount: result.Total,
	}).Error
	return
}

// ListChapters list the book's chapters
func (b *Book) ListChapters(bookID int, params *BookChapterQueryParams) (chapters []*model.Chapter, err error) {
	limit, _ := strconv.Atoi(params.Limit)
	offset, _ := strconv.Atoi(params.Offset)
	options := &model.QueryOptions{
		Limit:  limit,
		Offset: offset,
		Field:  params.Field,
		Order:  params.Order,
	}
	chapters = make([]*model.Chapter, 0)
	client := getClientByOptions(options, nil)
	err = client.Where(&model.Chapter{
		BookID: uint(bookID),
	}).Find(&chapters).Error
	return
}

// CountChapters count the book's chapters
func (b *Book) CountChapters(bookID int) (count int, err error) {
	err = getClient().
		Model(&model.Chapter{}).
		Where(&model.Chapter{
			BookID: uint(bookID),
		}).Count(&count).Error
	return
}

// incCount inc the count of book
func (b *Book) incCount(id int, field string) (err error) {
	book := &model.Book{}
	book.ID = uint(id)
	m := make(map[string]interface{})
	m[field] = gorm.Expr(field + " + 1")
	err = getClient().Model(book).Updates(m).Error
	return
}

// Like like the book
func (b *Book) Like(id int) (err error) {
	return b.incCount(id, bookLikeCount)
}

// View view the book
func (b *Book) View(id int) (err error) {
	return b.incCount(id, bookViewCount)
}

// UpdateCategories get category list
func (b *Book) UpdateCategories() (err error) {
	result := make([]pq.StringArray, 0)
	err = getClient().Model(&model.Book{}).Pluck(bookCategory, &result).Error
	if err != nil {
		return
	}
	categoriesInfo := make(map[string]int)
	for _, values := range result {
		for _, category := range values {
			categoriesInfo[category]++
		}
	}
	_, err = RedisSet(cs.CacheBookCategories, categoriesInfo, time.Hour)
	return
}
