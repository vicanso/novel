package xlog

import (
	"testing"
)

func TestGetLogger(t *testing.T) {
	if Logger() == nil {
		t.Fatalf("get logger fail")
	}

	if AccessLogger() == nil {
		t.Fatalf("create sugger logger fail")
	}

	if TrackerLogger() == nil {
		t.Fatalf("create tracker logger fail")
	}

	if UserLogger() == nil {
		t.Fatalf("create user logger fail")
	}
}
