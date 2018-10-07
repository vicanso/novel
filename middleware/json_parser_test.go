package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/xerror"
)

func TestNewJSONParser(t *testing.T) {
	fn := NewJSONParser(JSONParserConfig{
		Limit: 10 * 1024,
	})(func(c echo.Context) error {
		return nil
	})
	t.Run("pass get method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		if err != nil {
			t.Fatalf("get method, json parse fail, %v", err)
		}
	})

	t.Run("pass post body(not json)", func(t *testing.T) {
		body := []byte(`<xml></xml>`)
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/xml")
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		if err != nil {
			t.Fatalf("xml data, json parse fail, %v", err)
		}
	})
	t.Run("parse post body", func(t *testing.T) {
		body := []byte(`{"account": "vicanso"}`)
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		if err != nil {
			t.Fatalf("json parse fail, %v", err)
		}
		data := context.GetRequestBody(c)
		if !bytes.Equal(data, body) {
			t.Fatalf("json parse fail")
		}
	})

	t.Run("read post body fail", func(t *testing.T) {
		message := "read error"
		reader := xerror.NewErrorReadCloser(errors.New(message))
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/", reader)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		fmt.Println(err)
		if err == nil || err.Error() != message {
			t.Fatalf("read body should return error")
		}
	})

	t.Run("post nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(r, w)
		err := fn(c)
		if err != nil {
			t.Fatalf("post nil fail, %v", err)
		}
		data := context.GetRequestBody(c)
		if len(data) != 0 {
			t.Fatalf("the request body should be empty")
		}
	})

	t.Run("post body data parse over limit", func(t *testing.T) {
		limitFn := NewJSONParser(JSONParserConfig{
			Limit: 10,
		})(func(c echo.Context) error {
			return nil
		})
		body := []byte(`{"account": "vicanso"}`)
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(r, w)
		err := limitFn(c)
		if err == nil || err.(*xerror.HTTPError) != errJSONTooLarge {
			t.Fatalf("request post data should be too large")
		}
	})
}
