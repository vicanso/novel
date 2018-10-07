package util

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	now := Now()
	current := time.Now()
	_, z := current.Zone()
	size := 25
	// 如果刚好是0时区，那就不是25长度了(和utc一样)
	if z == 0 {
		size = 20
	}
	if len(now) != size {
		t.Fatalf("get now fail")
	}
}

func TestUTCNow(t *testing.T) {
	now := UTCNow()
	if len(now) != 20 {
		t.Fatalf("get utc now fail")
	}
}
