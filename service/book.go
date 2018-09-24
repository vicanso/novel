package service

import (
	"encoding/base64"
	"time"

	"github.com/vicanso/novel-spider/novel"
	"go.uber.org/zap"

	"github.com/vicanso/novel-spider/mq"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/util"
)

var (
	mqClient    *mq.MQ
	uploadURL   string
	uploadToken string
)

const (
	defaultNovelQueryLimit = 10
	bookCoverCategory      = "book-cover"
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

func init() {
	address := config.GetStringSlice("nsq.lookup.address")

	c := &mq.MQ{
		LookupAddress: address,
		Logger:        getLogger(),
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
	getLogger().Info("connect to nsq success")

	uploadURL = config.GetString("tiny.host") + config.GetString("tiny.upload")
	uploadToken = config.GetString("tiny.token")

}

func uploadCover(url string) (file string, err error) {
	if url == "" {
		return
	}
	buf, err := util.HTTPGet(url, nil)
	if err != nil {
		return
	}
	data := make(map[string]interface{})
	data["token"] = uploadToken
	data["category"] = bookCoverCategory
	data["fileType"] = "jpeg"
	data["maxAge"] = "720h"
	data["data"] = base64.StdEncoding.EncodeToString(buf)
	res, err := util.HTTPPost(uploadURL, data, nil)
	if err != nil {
		return
	}
	file = json.Get(res, "file").ToString()
	return
}

func initReceiveBasicInfoEvent(c *mq.MQ) (err error) {
	cb := func(info *novel.BasicInfo) {
		if info == nil {
			return
		}
		cover, _ := uploadCover(info.Cover)
		b := model.Book{
			Name:   info.Name,
			Author: info.Author,
			Brief:  info.Brief,
			Cover:  cover,
			Category: []string{
				info.Category,
			},
			Source:   info.Source,
			SourceID: info.SourceID,
		}
		err := model.AddBook(&b)
		if err != nil {
			getLogger().Error("add book fail",
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
	cb := func(c *novel.Chapter) {
		// 小于1000字的章节认为非正常数据
		wordCount := len(c.Content)
		if wordCount < 1000 {
			return
		}
		chapter := model.Chapter{
			Title:     c.Title,
			Content:   c.Content,
			WordCount: wordCount,
			Index:     c.Index,
		}
		model.AddBookChapter(c.Name, c.Author, &chapter)
	}
	err = c.SubReceiveChapter(cb)
	return
}

// AddBook add novel by category and id
func AddBook(category string, id int) (err error) {
	exists, _ := model.IsExistsBook(category, id)
	if exists {
		return
	}
	return mqClient.Pub(mq.TopicAddNovel, &novel.Source{
		Category: category,
		ID:       id,
	})
}

// AddBookChapter add chapter
func AddBookChapter(category string, id, chapterIndex int) (err error) {
	return mqClient.Pub(mq.TopicUpdateChapter, &novel.Source{
		Category:     category,
		ID:           id,
		ChapterIndex: chapterIndex,
	})
}
