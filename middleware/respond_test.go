package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/xerror"

	"github.com/labstack/echo"
)

func TestNewRespond(t *testing.T) {
	fn := NewRespond(RespondConfig{})(func(c echo.Context) error {
		return nil
	})
	t.Run("response no body", func(t *testing.T) {
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(nil, w)
		err := fn(c)
		if err != nil {
			t.Fatalf("respond middleware fail, %v", err)
		}
		if w.Code != http.StatusNoContent || len(w.Body.Bytes()) != 0 {
			t.Fatalf("response no body fail")
		}
	})

	t.Run("response string", func(t *testing.T) {
		text := "abcd"
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(nil, w)
		context.Res(c, text)
		err := fn(c)
		if err != nil {
			t.Fatalf("respond middleware fail, %v", err)
		}
		if w.Code != http.StatusOK || text != string(w.Body.Bytes()) {
			t.Fatalf("response string fail")
		}
	})

	t.Run("response string", func(t *testing.T) {
		text := "abcd"
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(nil, w)
		context.Res(c, text)
		err := fn(c)
		if err != nil {
			t.Fatalf("respond middleware fail, %v", err)
		}
		if w.Code != http.StatusOK || text != string(w.Body.Bytes()) {
			t.Fatalf("response string fail")
		}
	})

	t.Run("response bytes", func(t *testing.T) {
		buf := []byte("abcd")
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(nil, w)
		context.Res(c, buf)
		err := fn(c)
		if err != nil {
			t.Fatalf("respond middleware fail, %v", err)
		}
		header := w.HeaderMap
		if w.Code != http.StatusOK ||
			!bytes.Equal(buf, w.Body.Bytes()) ||
			header["Content-Type"][0] != "application/octet-stream" {
			t.Fatalf("response string fail")
		}
	})

	t.Run("response json", func(t *testing.T) {
		m := map[string]interface{}{
			"account": "vicanso",
			"age":     18,
			"vip":     true,
		}
		buf, _ := json.Marshal(m)
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(nil, w)
		context.Res(c, m)
		err := fn(c)
		if err != nil {
			t.Fatalf("respond middleware fail, %v", err)
		}
		if w.Code != http.StatusOK ||
			!bytes.Equal(buf, w.Body.Bytes()) ||
			w.Header()["Content-Type"][0] != "application/json; charset=UTF-8" {
			t.Fatalf("response json fail")
		}
	})

	t.Run("response error", func(t *testing.T) {
		errFn := NewRespond(RespondConfig{})(func(c echo.Context) error {
			return errors.New("abcd")
		})
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(nil, w)
		err := errFn(c)
		if err != nil {
			t.Fatalf("respond middleware fail, %v", err)
		}
		if w.Code != http.StatusInternalServerError ||
			string(w.Body.Bytes()) != `{"statusCode":500,"exception":true,"message":"abcd","category":"exception"}` {
			t.Fatalf("response error fail")
		}
	})

	t.Run("response http error", func(t *testing.T) {
		errFn := NewRespond(RespondConfig{})(func(c echo.Context) error {
			he := &xerror.HTTPError{
				StatusCode: 401,
				Category:   "a",
				Exception:  false,
				Message:    "abcd",
			}
			return he
		})
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(nil, w)
		err := errFn(c)
		if err != nil {
			t.Fatalf("respond middleware fail, %v", err)
		}
		if w.Code != 401 ||
			string(w.Body.Bytes()) != `{"statusCode":401,"message":"abcd","category":"a"}` {
			t.Fatalf("response http error fail")
		}
	})
}
