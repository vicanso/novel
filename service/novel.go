package service

import (
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/vicanso/novel-spider/novel"

	"github.com/vicanso/novel-spider/mq"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/utils"
)

var (
	mqClient *mq.MQ
)

const (
	defaultNovelQueryLimit = 10
)

type (
	// NovelQueryParams novel query params
	NovelQueryParams struct {
		Limit     string `json:"limit,omitempty" valid:"in(1|10|20|50|100),optional"`
		Offset    string `json:"offset,omitempty" valid:"int,optional"`
		Field     string `json:"field,omitempty" valid:"runelength(0|1000),optional"`
		Order     string `json:"order,omitempty" valid:"runelength(0|100),optional"`
		Name      string `json:"name,omitempty" valid:"runelength(0|20),optional"`
		Author    string `json:"author,omitempty" valid:"runelength(0|20),optional"`
		Status    string `json:"status,omitempty" valid:"in(1|2|3),optional"`
		Category  string `json:"category,omitempty" valid:"runelength(0|10),optional"`
		UpdatedAt string `json:"updatedAt,omitempty" valid:"runelength(0|100),optional"`
		CreatedAt string `json:"createdAt,omitempty" valid:"runelength(0|100),optional"`
	}
)

func initReceiveBasicInfoEvent(c *mq.MQ) (err error) {
	cb := func(info *novel.BasicInfo) {
		err := model.AddBook(info)
		if err != nil {
			utils.GetLogger().Error("add book fail",
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
	cb := func(chapter *novel.Chapter) {
		fmt.Println(chapter)
	}
	err = c.SubReceiveChapter(cb)
	return
}

func init() {
	address := config.GetStringSlice("nsq.lookup.address")

	c := &mq.MQ{
		LookupAddress: address,
	}
	c.FreshNodes()
	go c.TimedFreshNodes(time.Second * 60)

	if err := initReceiveBasicInfoEvent(c); err != nil {
		panic(err)
	}
	if err := initReceiveChapterEvent(c); err != nil {
		panic(err)
	}
	mqClient = c
	utils.GetLogger().Info("connect to nsq success")
}

func getBookQueryConditions(params *NovelQueryParams) *model.BookQueryConditions {
	return &model.BookQueryConditions{
		Name:      params.Name,
		Author:    params.Author,
		Status:    params.Status,
		Category:  params.Category,
		UpdatedAt: params.UpdatedAt,
		CreatedAt: params.CreatedAt,
	}
}

// AddNovel add novel by category and id
func AddNovel(category string, id int) (err error) {
	exists, _ := model.IsExistsSource(category, id)
	if exists {
		return
	}
	return mqClient.Pub(mq.TopicAddNovel, &novel.Source{
		Category: category,
		ID:       id,
	})
}

// ListNovel list the novel
func ListNovel(params *NovelQueryParams) (books []*model.Book, err error) {
	options := &model.QueryOptions{
		Field: params.Field,
		Order: params.Order,
		Limit: defaultNovelQueryLimit,
	}
	if params.Limit != "" {
		options.Limit, err = strconv.Atoi(params.Limit)
		if err != nil {
			return
		}
	}
	if params.Offset != "" {
		options.Offset, err = strconv.Atoi(params.Offset)
		if err != nil {
			return
		}
	}

	conditions := getBookQueryConditions(params)
	books, err = model.ListBook(conditions, options)
	return
}

// CountNovel count novel
func CountNovel(params *NovelQueryParams) (count int, err error) {
	conditions := getBookQueryConditions(params)
	count, err = model.CountBook(conditions)
	return
}
