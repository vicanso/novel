package util

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

func TestSha256(t *testing.T) {
	if Sha256("vicanso") != "WLNYrrG6ofufjI5HDURypQhIRfsADsWAjZkSXT1ilkQ=" {
		t.Fatalf("sha256 fail")
	}
}

func TestContainsString(t *testing.T) {
	if !ContainsString([]string{"a", "b"}, "a") {
		t.Fatalf("contains fail")
	}
	if ContainsString([]string{"a", "b"}, "c") {
		t.Fatalf("contains fail")
	}
}
