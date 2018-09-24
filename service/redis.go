package service

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/vicanso/novel/config"
)

var (
	redisClient   *redis.Client
	redisOkResult = "OK"
	redisNoop     = func() error {
		return nil
	}
)

type (
	// RedisDoneFn redis done function
	RedisDoneFn func() error
)

func init() {
	uri := config.GetString("redis")
	if uri != "" {
		redisClient = newRedisClient(uri)
	}
}

// newRedisClient new client
func newRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
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
