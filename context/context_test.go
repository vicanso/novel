package context

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/service"
)

func TestSetHeader(t *testing.T) {
	e := echo.New()
	res := httptest.NewRecorder()
	c := e.NewContext(nil, res)
	k := "X-Token"
	v := "b"
	SetHeader(c, k, v)
	if c.Response().Header()[k][0] != v {
		t.Fatalf("set header fail")
	}
}

func TestSetCache(t *testing.T) {
	e := echo.New()
	res := httptest.NewRecorder()
	c := e.NewContext(nil, res)
	t.Run("no cache", func(t *testing.T) {
		SetNoCache(c)
		if c.Response().Header()["Cache-Control"][0] != "no-cache, max-age=0" {
			t.Fatalf("set no cache fail")
		}
	})
	t.Run("set max age", func(t *testing.T) {
		SetCache(c, "10s")
		if c.Response().Header()["Cache-Control"][0] != "public, max-age=10" {
			t.Fatalf("set max age fail")
		}
	})
	t.Run("set s-max age", func(t *testing.T) {
		SetCacheWithSMaxAge(c, "10s", "1s")
		if c.Response().Header()["Cache-Control"][0] != "public, max-age=10, s-maxage=1" {
			t.Fatalf("set s-max age fail")
		}
	})
	t.Run("set private cache", func(t *testing.T) {
		SetPrivateCache(c, "10s")
		if c.Response().Header()["Cache-Control"][0] != "private, max-age=10" {
			t.Fatalf("set private max age fail")
		}
	})
}

func TestGetTrackID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
	trackID := "random-string"

	e := echo.New()
	c := e.NewContext(r, nil)
	if GetTrackID(c) != "" {
		t.Fatalf("get nil track fail")
	}

	r.AddCookie(&http.Cookie{
		Name:  config.GetTrackKey(),
		Value: trackID,
	})
	if GetTrackID(c) != trackID {
		t.Fatalf("get track id fail")
	}
}

func TestStatus(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)
	if GetStatus(c) != 204 {
		t.Fatalf("the original status should be 204")
	}
	Res(c, []byte("abcd"))
	if GetStatus(c) != http.StatusOK {
		t.Fatalf("the status should be 200")
	}
	SetStatus(c, http.StatusCreated)
	if GetStatus(c) != http.StatusCreated {
		t.Fatalf("set the status fail")
	}
}

func TestContentType(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)
	if GetContentType(c) != "" {
		t.Fatalf("the content type should be nil")
	}
	SetContentType(c, echo.MIMEApplicationJSON)
	if GetContentType(c) != echo.MIMEApplicationJSON {
		t.Fatalf("set content type fail")
	}
}

func TestRes(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)
	data := map[string]string{
		"name": "tree.xie",
	}
	Res(c, data)
	v := GetBody(c)
	if v == nil || v.(map[string]string)["name"] != "tree.xie" {
		t.Fatalf("get and set response body fail")
	}

	ResCreated(c, data)
	if GetStatus(c) != http.StatusCreated {
		t.Fatalf("res create fail")
	}
}

func TestUserSession(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)

	if GetUserSession(c) != nil {
		t.Fatalf("get session should be nil before set")
	}
	us := &service.UserSession{}
	SetUserSession(c, us)
	if GetUserSession(c) != us {
		t.Fatalf("get/set session fail")
	}
}

func TestRequestBody(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)

	buf := []byte("abcd")
	if GetRequestBody(c) != nil {
		t.Fatalf("get request body should be nil before set")
	}
	SetRequestBody(c, buf)
	if !bytes.Equal(buf, GetRequestBody(c)) {
		t.Fatalf("get request body fail")
	}
}

func TestRequestQuery(t *testing.T) {
	e := echo.New()
	t.Run("no query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		c := e.NewContext(req, nil)
		if GetRequestQuery(c) != nil {
			t.Fatalf("get nil query fail")
		}
	})
	t.Run("query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/?a=1&b=2", nil)
		c := e.NewContext(req, nil)
		GetRequestQuery(c)
		q := GetRequestQuery(c)
		if q["a"] != "1" || q["b"] != "2" {
			t.Fatalf("get request query fail")
		}
	})
}

func TestGetRequestParams(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)
	params := GetRequestParams(c)
	if len(params) != 0 {
		t.Fatalf("nil params fail")
	}
	c.SetParamNames("a", "b")
	c.SetParamValues("1", "2")
	params = GetRequestParams(c)
	if params["a"] != "1" && params["b"] != "2" {
		t.Fatalf("get params fail")
	}

}

func TestGetRequestID(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)
	id := GetRequestID(c)
	if id == "" {
		t.Fatalf("get request id fail")
	}
	if GetRequestID(c) != id {
		t.Fatalf("get request id(twice) fail")
	}
}

func TestGetLogger(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)
	GetLogger(c)
	if GetLogger(c) == nil {
		t.Fatalf("get context logger fail")
	}
}
