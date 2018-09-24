package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	"github.com/vicanso/novel/util"

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
	// likePrefix like query prefix
	likePrefix = '~'
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

func initModels() {
	db.AutoMigrate(&User{}).
		AutoMigrate(&UserLogin{}).
		AutoMigrate(&Book{}).
		AutoMigrate(&Chapter{})
}

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
	util.GetLogger().Info("connect to postgres success",
		zap.String("uri", str),
	)
}

// GetClient get db client
func GetClient() *gorm.DB {
	return db
}

func getClientByOptions(options *QueryOptions, client *gorm.DB) *gorm.DB {
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

// enhanceWhere 增强的enhance where查询
func enhanceWhere(client *gorm.DB, key, value string) *gorm.DB {
	if key == "" || value == "" {
		return client
	}
	if value[0] == likePrefix {
		like := key + " LIKE ?"
		value = "%" + value[1:] + "%"
		client = client.Where(like, value)
	} else {
		client = client.Where(key+" = ?", value)
	}
	return client
}

// enhanceRangeWhere 增强的区间查询
func enhanceRangeWhere(client *gorm.DB, key, value string) *gorm.DB {
	if key == "" || value == "" {
		return client
	}
	arr := strings.Split(value, "|")
	start := arr[0]
	end := ""
	if len(arr) > 1 {
		end = arr[1]
	}
	if start != "" {
		client = client.Where(key+" >= ?", start)
	}
	if end != "" {
		client = client.Where(key+" <= ?", end)
	}
	return client
}
