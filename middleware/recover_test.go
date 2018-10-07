package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

func TestNewRecover(t *testing.T) {
	e := echo.New()
	t.Run("panic error", func(t *testing.T) {
		fn := NewRecover(RecoverConfig{})(func(c echo.Context) error {
			panic("abcd")
		})
		req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		err := fn(c)
		if err != nil {
			t.Fatalf("recover middleware fail, %v", err)
		}
	})
}
