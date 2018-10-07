package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/cs"
)

func TestNewEntry(t *testing.T) {
	e := echo.New()
	res := httptest.NewRecorder()
	c := e.NewContext(nil, res)
	k := "X-Token"
	v := "a"
	fn := NewEntry(EntryConfig{
		Header: map[string]string{
			k: v,
		},
	})(func(c echo.Context) error {
		return nil
	})
	err := fn(c)
	if err != nil {
		t.Fatalf("entry middleware fail, %v", err)
	}
	header := c.Response().Header()
	if header[k][0] != v {
		t.Fatalf("set header fail")
	}
	if header[cs.HeaderCacheControl][0] != cs.HeaderNoCache {
		t.Fatalf("set no cache fail")
	}
}
