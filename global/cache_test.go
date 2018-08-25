package global

import (
	"testing"

	"github.com/vicanso/novel/utils"
)

func TestCache(t *testing.T) {
	key := utils.RandomString(8)
	value := "a"
	Store(key, value)
	v, ok := Load(key)
	if !ok || v.(string) != value {
		t.Fatalf("store and load cache fail")
	}
	_, loaded := LoadOrStore(key, "b")
	if !loaded {
		t.Fatalf("load or store should loaded while data exists")
	}

	key = utils.RandomString(8)
	v, loaded = LoadOrStore(key, "b")
	if loaded {
		t.Fatalf("load or store should not be loaded while data not exists")
	}
	if v.(string) != "b" {
		t.Fatalf("load or store fail")
	}

}
