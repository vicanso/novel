package utils

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	size := 8
	v := RandomString(size)
	if len(v) != size {
		t.Fatalf("create random string fail")
	}
}

func TestUlid(t *testing.T) {
	id := GenUlid()
	if len(id) != 26 {
		t.Fatalf("create ulid string fail")
	}
}
