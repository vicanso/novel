package global

import (
	"sync"
)

var (
	m = &sync.Map{}
)

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
