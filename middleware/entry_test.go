package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/novel/utils"
)

func TestNewEntry(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
	w := httptest.NewRecorder()
	fn := NewEntry()
	ctx := utils.NewContext(w, r)
	fn(ctx)
	logger := utils.GetLogger()
	if logger == nil {
		t.Fatalf("entry middle should create a user logger")
	}
}
