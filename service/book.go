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
	bookCoverCategory   = "book-cover"
	bookViewCount       = "view_count"
	bookLatestViewCount = "latest_view_count"
	bookLikeCount       = "like_count"
	bookLatestLikeCount = "latest_like_count"
	bookWordCount       = "word_count"
	bookCategory        = "category"
	bookAuthor          = "author"
	bookName            = "name"
	bookStatus          = "status"
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
		Author   string `json:"author,omitempty" valid:"runelength(2|20),optional"`
		Status   string `json:"status,omitempty" valid:"in(0|1|2),optional"`
	}
	// BookUpdateParams params for update
	BookUpdateParams struct {
		Name        string `json:"name,omitempty" valid:"runelength(1|20),optional"`
		Brief       string `json:"brief,omitempty" valid:"runelength(5|2000),optional"`
		Status      int    `json:"status,omitempty" valid:"xIntIn(1|2),optional"`
		Category    string `json:"category,omitempty" valid:"runelength(2|32),optional"`
		SourceCover string `json:"sourceCover" valid:"runelength(10|200),optional"`
	}
	// BookChapterQueryParams params for the query
	BookChapterQueryParams struct {
		Limit  string `json:"limit,omitempty" valid:"range(1|100)"`
		Offset string `json:"offset,omitempty" valid:"numeric"`
		Field  string `json:"field,omitempty" valid:"runelength(2|64)"`
		Order  string `json:"order,omitempty" valid:"runelength(2|32)"`
	}
	// BookFavorite book favorite
	BookFavorite struct {
		ID        uint       `json:"id,omitempty"`
		Name      string     `json:"name,omitempty"`
		Author    string     `json:"author,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		// 收藏时间
		CreatedAt      *time.Time     `json:"createdAt,omitempty"`
		LatestChapter  *model.Chapter `json:"latestChapter,omitempty"`
		ReadingChapter *model.Chapter `json:"readingChapter,omitempty"`
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
				Source:   c.Source,
				SourceID: c.SourceID,
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
func (b *Book) UpdateChapters(id int) (err error) {
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
		Order: "index",
	}, nil).
		Where(&model.Chapter{
			BookID: book.ID,
		}).
		Find(&chapters).
		Error
	if err != nil {
		return
	}
	index := 0
	for _, chapter := range chapters {
		if chapter.Index > index {
			// 对缺失的章节更新
			for i := index; i < chapter.Index; i++ {
				b.AddChapter(book.Source, book.SourceID, i)
			}
		}
		index = chapter.Index + 1
	}
	latestChapter := 0
	chapterCount := len(chapters)
	if chapterCount != 0 {
		latestChapter = chapters[chapterCount-1].Index + 1
	}
	// 更新最新章节
	err = b.AddLatestChapter(book.Source, book.SourceID, latestChapter)
	return
}

// AddChapter add chapter
func (b *Book) AddChapter(category string, id, index int) (err error) {
	return mqClient.Pub(mq.TopicUpdateChapter, &novel.Source{
		Category:   category,
		ID:         id,
		Chapter:    index,
		UpdateType: novel.UpdateTypeCurrent,
	})
}

// AddLatestChapter add latest chapter
func (b *Book) AddLatestChapter(category string, id, index int) (err error) {
	return mqClient.Pub(mq.TopicUpdateChapter, &novel.Source{
		Category:   category,
		ID:         id,
		Chapter:    index,
		UpdateType: novel.UpdateTypeLatest,
	})
}

// UpdateCover update cover
func (b *Book) UpdateCover(id int) (file string, err error) {
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
	file = json.Get(res, "file").ToString()
	err = getClient().Model(book).Update(&model.Book{
		Cover: file,
	}).Error
	return
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
	author := params.Author
	if author != "" {
		sql = append(sql, bookAuthor+" = ?")
		args = append(args, author)
	}
	query = strings.Join(sql, " AND ")
	return
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
	if params.Name != "" {
		updateInfo.Name = params.Name
	}
	if params.SourceCover != "" {
		updateInfo.SourceCover = params.SourceCover
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
func (b *Book) incCount(id int, fields []string) (err error) {
	book := &model.Book{}
	book.ID = uint(id)
	m := make(map[string]interface{})
	for _, field := range fields {
		m[field] = gorm.Expr(field + " + 1")
	}
	err = getClient().Model(book).Updates(m).Error
	return
}

// Like like the book
func (b *Book) Like(id int) (err error) {
	return b.incCount(id, []string{
		bookLikeCount,
		bookLatestLikeCount,
	})
}

// View view the book
func (b *Book) View(id int) (err error) {
	return b.incCount(id, []string{
		bookViewCount,
		bookLatestViewCount,
	})
}

// ResetLatestCount reset latest count
func (b *Book) ResetLatestCount() (err error) {
	err = getClient().
		Model(&model.Book{}).
		Where(&model.Book{
			Status: model.BookStatusPassed,
		}).
		Updates(map[string]interface{}{
			bookLatestLikeCount: 0,
			bookLatestViewCount: 0,
		}).Error
	return
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

// GetLatestChapter get latest chapter
func (b *Book) GetLatestChapter(bookID int, field string) (chapter *model.Chapter, err error) {
	chapters, err := b.ListChapters(bookID, &BookChapterQueryParams{
		Limit: "1",
		Order: "-index",
		Field: field,
	})
	if err != nil || len(chapters) == 0 {
		return
	}
	chapter = chapters[0]
	return
}

// AddFav add fav book
func (b *Book) AddFav(account string, bookID int) (err error) {
	result := &model.Book{}
	result.ID = uint(bookID)
	err = getClient().First(result).Error
	if err != nil {
		return
	}
	fav := &model.Favorite{
		Account: account,
		BookID:  uint(bookID),
	}
	err = getClient().Create(fav).Error
	return
}

// RemoveFav remove fav book
func (b *Book) RemoveFav(account string, bookID int) (err error) {
	err = getClient().Unscoped().Delete(&model.Favorite{
		Account: account,
		BookID:  uint(bookID),
	}).Error
	return
}

// ListFav list fav books
func (b *Book) ListFav(account string) (bookFavs []*BookFavorite, err error) {
	favs := make([]*model.Favorite, 0)
	err = getClient().Where(&model.Favorite{
		Account: account,
	}).Find(&favs).Error
	if err != nil {
		return
	}

	books := make([]*model.Book, 0)
	ids := []uint{}
	for _, item := range favs {
		ids = append(ids, item.BookID)
	}
	err = getClientByOptions(&model.QueryOptions{
		Field: "id,name,author,updatedAt",
	}, nil).Where("id in (?)", ids).
		Find(&books).Error
	if err != nil {
		return
	}
	bookFavs = make([]*BookFavorite, len(books))
	waitChans := make(chan bool, 5)
	for index, item := range books {
		bookFav := &BookFavorite{
			ID:        item.ID,
			Name:      item.Name,
			Author:    item.Author,
			UpdatedAt: item.UpdatedAt,
		}
		go func() {
			var foundFav *model.Favorite
			for _, fav := range favs {
				if fav.BookID == item.ID {
					bookFav.CreatedAt = fav.CreatedAt
					foundFav = fav
				}
			}
			chapterFields := "title,index,updatedAt"
			latestChapter, _ := b.GetLatestChapter(int(bookFav.ID), chapterFields)
			if latestChapter != nil {
				bookFav.LatestChapter = latestChapter
			}
			if foundFav != nil {
				bookFav.ReadingChapter = &model.Chapter{}
				getClientByOptions(&model.QueryOptions{
					Field: chapterFields,
				}, nil).Where(&model.Chapter{
					BookID: bookFav.ID,
					Index:  foundFav.ReadingChapter,
				}).First(bookFav.ReadingChapter)
			}
			waitChans <- true
		}()
		<-waitChans
		bookFavs[index] = bookFav
	}
	return
}

// UpdateFav update the fav info
func (b *Book) UpdateFav(account string, bookID, readingChapter int) (err error) {
	err = getClient().
		Model(&model.Favorite{}).
		Where(&model.Favorite{
			Account: account,
			BookID:  uint(bookID),
		}).
		Updates(&model.Favorite{
			ReadingChapter: readingChapter,
		}).Error
	return
}
