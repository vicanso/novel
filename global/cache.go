package global

// 全局使用的缓存配置，可用于缓存一些全局信息

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/golang-lru"
)

const (
	defaultCacheSize   = 1024 * 10
	connectingCountKey = "connecting"
	routeInfosKey      = "routeInfo"
)

type (
	// RouteCounter route counter
	RouteCounter struct {
		CreatedAt string
		Counts    map[string]*uint32
	}
)

var (
	m        = &sync.Map{}
	lruCache *lru.Cache
	// routeCounter the route counter info
	routeCounter = &RouteCounter{}
)

func init() {
	l, err := lru.New(defaultCacheSize)
	if err != nil {
		panic(err)
	}
	lruCache = l
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

// SaveConnectingCount save the current connecting count
func SaveConnectingCount(v uint32) {
	m.Store(connectingCountKey, v)
}

// GetConnectingCount get the current connecting count
func GetConnectingCount() (connectingCount uint32) {
	v, ok := m.Load(connectingCountKey)
	if !ok || v == nil {
		return 0
	}
	return v.(uint32)
}

// SaveRouteInfos save route infos
func SaveRouteInfos(v []map[string]string) {
	m.Store(routeInfosKey, v)
}

// GetRouteInfos get route infos
func GetRouteInfos() (routeInfo []map[string]string) {
	v, ok := m.Load(routeInfosKey)
	if !ok || v == nil {
		return nil
	}
	return v.([]map[string]string)
}

// Load get data from cache
func Load(key interface{}) (interface{}, bool) {
	return m.Load(key)
}

// Store store data to cache
func Store(key, value interface{}) {
	m.Store(key, value)
}

// LoadOrStore load the data from cache, if not exists, store it
func LoadOrStore(key, value interface{}) (interface{}, bool) {
	return m.LoadOrStore(key, value)
}

// NewLRU new a lru cache
func NewLRU(size int) (*lru.Cache, error) {
	return lru.New(size)
}

// Add add value to lru cache（default cache）
func Add(key, value interface{}) (evicted bool) {
	return lruCache.Add(key, value)
}

// Get get the value from lru cache
func Get(key interface{}) (value interface{}, found bool) {
	return lruCache.Get(key)
}

// Remove remove the key from lru cache
func Remove(key interface{}) {
	lruCache.Remove(key)
}

// InitRouteCounter init route counter
func InitRouteCounter(routeInfos []map[string]string) {
	routeCounter.CreatedAt = now()
	routeCounter.Counts = make(map[string]*uint32)
	counts := routeCounter.Counts
	for _, info := range routeInfos {
		key := info["method"] + " " + info["path"]
		var v uint32
		counts[key] = &v
	}
}

// AddRouteCount add the route's count
func AddRouteCount(method, path string) {
	if method == "" || path == "" {
		return
	}
	key := method + " " + path
	v := routeCounter.Counts[key]
	if v == nil {
		return
	}
	atomic.AddUint32(v, 1)
}

// ResetRouteCount reset the route count
func ResetRouteCount() {
	for _, v := range routeCounter.Counts {
		atomic.StoreUint32(v, 0)
	}
	routeCounter.CreatedAt = now()
}

// GetRouteCount get the route count
func GetRouteCount() map[string]interface{} {
	m := make(map[string]uint32)
	for k, v := range routeCounter.Counts {
		m[k] = *v
	}
	data := make(map[string]interface{})
	data["createdAt"] = routeCounter.CreatedAt
	data["counts"] = m
	return data
}
