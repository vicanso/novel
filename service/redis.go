package service

import (
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/xlog"
	"go.uber.org/zap"
)

var (
	redisClient   *redis.Client
	redisOkResult = "OK"
	redisNoop     = func() error {
		return nil
	}
	// ErrRedisNil redis nil error
	ErrRedisNil = redis.Nil
)

type (
	// RedisDoneFn redis done function
	RedisDoneFn func() error
)

func init() {
	uri := config.GetString("redis")
	if uri != "" {
		c, err := newRedisClient(uri)
		if err != nil {
			panic(err)
		}
		redisClient = c
		_, err = redisClient.Ping().Result()
		logger := xlog.Logger()
		mask := regexp.MustCompile(`redis://:(\S+)\@`)
		str := mask.ReplaceAllString(uri, "redis://:***@")
		if err != nil {
			logger.Error("redis ping fail",
				zap.String("uri", str),
				zap.Error(err),
			)
		} else {
			logger.Info("redis ping success",
				zap.String("uri", str),
			)
		}
	}
}

// newRedisClient new client
func newRedisClient(uri string) (client *redis.Client, err error) {
	info, err := url.Parse(uri)
	if err != nil {
		return
	}
	opts := &redis.Options{
		Addr: info.Host,
	}
	db := info.Query().Get("db")
	if db != "" {
		opts.DB, _ = strconv.Atoi(db)
	}
	opts.Password, _ = info.User.Password()
	client = redis.NewClient(opts)
	return
}

// GetRedisClient get redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// Lock lock the key for ttl seconds
func Lock(key string, ttl time.Duration) (bool, error) {
	return redisClient.SetNX(key, true, ttl).Result()
}

// LockWithDone lock the key for ttl, and with done function
func LockWithDone(key string, ttl time.Duration) (bool, RedisDoneFn, error) {
	success, err := Lock(key, ttl)
	// 如果lock失败，则返回no op的done function
	if err != nil || !success {
		return false, redisNoop, err
	}
	done := func() error {
		_, err := redisClient.Del(key).Result()
		return err
	}
	return true, done, nil
}

// RedisSet the cache with ttl
func RedisSet(key string, v interface{}, ttl time.Duration) (ok bool, err error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return
	}
	result, err := redisClient.Set(key, buf, ttl).Result()
	if err != nil {
		return
	}
	ok = result == redisOkResult
	return
}

// RedisGet get the cache to v
func RedisGet(key string, v interface{}) (ok bool, err error) {
	buf, err := redisClient.Get(key).Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, v)
	if err != nil {
		return
	}
	ok = true
	return
}

// IsRedisNil check the error is redis nil
func IsRedisNil(err error) bool {
	if err == nil {
		return false
	}
	// 无法获取到RedisError对象
	return err.Error() == redis.Nil.Error()
}
