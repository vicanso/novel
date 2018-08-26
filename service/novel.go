package service

import (
	"fmt"
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
