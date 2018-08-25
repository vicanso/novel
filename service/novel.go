package service

import (
	"fmt"
	"time"

	"github.com/vicanso/novel-spider/novel"

	"github.com/vicanso/novel-spider/mq"
	"github.com/vicanso/novel/config"
)

var (
	mqClient *mq.MQ
)

func initReceiveBasicInfoEvent(c *mq.MQ) (err error) {
	cb := func(info *novel.BasicInfo) {
		fmt.Println(info)
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
}

// AddNovel add novel by category and id
func AddNovel(category string, id int) (err error) {
	return mqClient.Pub(mq.TopicAddNovel, &novel.Source{
		Category: category,
		ID:       id,
	})
}
