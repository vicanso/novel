package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/vicanso/novel/xlog"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"

	// for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var (
	db          *gorm.DB
	matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
)

const (
	defaultPoolSize = 100
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
	// QueryOptions query options
	QueryOptions struct {
		Limit  int
		Offset int
		Field  string
		Order  string
	}
)

// Init 初始化
func init() {
	uri := viper.GetString("db.uri")
	client, err := gorm.Open("postgres", uri)
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

	mask := regexp.MustCompile(`postgres://(\S+):(\S+)\@`)
	str := mask.ReplaceAllString(uri, "postgres://***:***@")
	xlog.Logger().Info("connect to postgres success",
		zap.String("uri", str),
		zap.Int("poolSize", poolSize),
	)
}

// init all models
func initModels() {
	db.AutoMigrate(&User{}).
		AutoMigrate(&UserLogin{}).
		AutoMigrate(&Book{}).
		AutoMigrate(&Chapter{})
}

// toSnakeCase convert string to snake case
func toSnakeCase(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// convertOrder convert the order string
func convertOrder(str string) string {
	arr := strings.Split(str, ",")
	result := make([]string, len(arr))
	for i, v := range arr {
		sort := "asc"
		if v[0] == '-' {
			sort = "desc"
			v = v[1:]
		}
		result[i] = toSnakeCase(v) + " " + sort
	}
	return strings.Join(result, ",")
}

// GetClient get db client
func GetClient() *gorm.DB {
	return db
}

// GetClientByOptions get db client by options
func GetClientByOptions(options *QueryOptions, client *gorm.DB) *gorm.DB {
	if client == nil {
		client = GetClient()
	}
	if options == nil {
		return client
	}
	if options.Field != "" {
		client = client.Select(toSnakeCase(options.Field))
	}
	if options.Order != "" {
		client = client.Order(convertOrder(options.Order))
	}
	if options.Limit > 0 {
		client = client.Limit(options.Limit)
	}
	if options.Offset > 0 {
		client = client.Offset(options.Offset)
	}
	return client
}
