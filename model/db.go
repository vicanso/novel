package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/vicanso/novel/utils"

	// for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var db *gorm.DB

const (
	defaultPoolSize = 100
)

const (
	// StatusPadding 待审核
	StatusPadding = iota + 1
	// StatusRejected 拒绝
	StatusRejected
	// StatusPassed 已通过
	StatusPassed
)

type (
	// BaseModel 基础的model定义
	BaseModel struct {
		// 唯一id
		ID uint `gorm:"primary_key" json:"id,omitempty"`
		// 创建时间
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		// 更新时间
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		// 删除时间
		DeletedAt *time.Time `json:"deletedAt,omitempty" sql:"index"`
	}
)

func initModels() {
	db.AutoMigrate(&Book{})
}

// Init 初始化
func init() {
	client, err := gorm.Open("postgres", viper.GetString("db.uri"))
	if err != nil {
		panic(fmt.Errorf("Fatal open postgres: %s", err))
	}
	poolSize := viper.GetInt("db.poolSize")
	if poolSize == 0 {
		poolSize = defaultPoolSize
	}
	client.DB().SetMaxOpenConns(poolSize)
	db = client
	initModels()
	utils.GetLogger().Info("connect to postgres success")
}

// GetClient 获取db client
func GetClient() *gorm.DB {
	return db
}
