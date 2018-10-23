package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/novel/asset"
	"github.com/vicanso/novel/xerror"

	"github.com/labstack/echo"
)

func TestServe(t *testing.T) {
	e := echo.New()
	adminAsset := asset.GetAdminAsset()
	t.Run("get gzip data", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := e.NewContext(nil, w)
		c.SetParamNames("*")
		c.SetParamValues("index.html")
		err := createServe(adminAsset)(c)
		if err != nil {
			t.Fatalf("static serve fail, %v", err)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		c := e.NewContext(nil, nil)
		c.SetParamNames("*")
		c.SetParamValues("a.html")
		err := createServe(adminAsset)(c).(*xerror.HTTPError)
		if err.Category != staticErrCategory ||
			err.StatusCode != http.StatusNotFound {
			t.Fatalf("file not found error is invalid")
		}
	})

}
