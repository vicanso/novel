package controller

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vicanso/novel/xerror"

	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/session"

	"github.com/labstack/echo"
)

func TestUserCtrl(t *testing.T) {
	e := echo.New()
	ctrl := UserCtrl{}
	cookies := []string{}
	account := util.RandomString(10)
	password := util.Sha256(config.GetString("app") + "123456")
	t.Run("get info", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/users/v1/me", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		sess := &session.Session{}
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		err := ctrl.getInfo(c)
		if err != nil {
			t.Fatalf("get info fail, %v", err)
		}
		cookies = c.Response().Header()["Set-Cookie"]
		userInfo := context.GetBody(c).(*UserInfoRes)
		if !userInfo.Anonymous {
			t.Fatalf("user info should be anonymous")
		}

		if userInfo.Date == "" {
			t.Fatalf("user info's date should not be empty")
		}
	})

	t.Run("get avatar", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := e.NewContext(nil, w)
		err := ctrl.getAvatar(c)
		if err != nil {
			t.Fatalf("get avatar fail, %v", err)
		}
		buf := context.GetBody(c).([]byte)
		if base64.StdEncoding.EncodeToString(buf) != avatar {
			t.Fatalf("the avatar content is wrong")
		}
	})

	t.Run("get login token", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/users/v1/me/token", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		sess := &session.Session{}
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}
		err := ctrl.getLoginToken(c)
		if err != nil {
			t.Fatalf("get login token fail, %v", err)
		}
		data := context.GetBody(c).(map[string]string)
		if len(data["token"]) != loginTokenLength {
			t.Fatalf("get login token fail")
		}
	})

	t.Run("register", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		context.SetRequestBody(c, []byte("{}"))
		err := ctrl.register(c)
		if err == nil {
			t.Fatalf("validte should fail")
		}
		he := err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategoryValidte {
			t.Fatalf("should be validate fail error")
		}

		m := map[string]string{
			"account":  account,
			"password": password,
		}
		buf, _ := json.Marshal(m)
		context.SetRequestBody(c, buf)
		err = ctrl.register(c)
		if err != nil {
			t.Fatalf("register fail, %v", err)
		}
		if context.GetStatus(c) != 201 {
			t.Fatalf("status should be 201")
		}
	})

	t.Run("login params invalid", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		context.SetRequestBody(c, []byte("{}"))
		err := ctrl.doLogin(c)
		if err == nil {
			t.Fatalf("validte should fail")
		}
		he := err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategoryValidte {
			t.Fatalf("should be validate fail error")
		}
	})

	t.Run("login token is nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		c := e.NewContext(r, w)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		context.SetRequestBody(c, []byte(`{
			"account": "vicanso",
			"password": "12341234"
		}`))
		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}
		err := ctrl.doLogin(c)
		if err == nil {
			t.Fatalf("should be fail when login token nil")
		}
		he := err.(*xerror.HTTPError)
		if he != errLoginTokenNil {
			t.Fatalf("the error should be login token nil")
		}
	})

	t.Run("login account/password is wrong", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		c := e.NewContext(r, w)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := service.NewUserSession(sess)
		context.SetUserSession(c, us)
		context.SetRequestBody(c, []byte(`{
			"account": "xxxxxxxx",
			"password": "12341234"
		}`))
		us.SetLoginToken(util.RandomString(loginTokenLength))

		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}

		err := ctrl.doLogin(c)
		if err == nil {
			t.Fatalf("should be fail when account is not exists")
		}
		he := err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategoryUser {
			t.Fatalf("the error category should be user")
		}

		context.SetRequestBody(c, []byte(`{
			"account": "`+account+`",
			"password": "12341234"
		}`))

		err = ctrl.doLogin(c)
		if err == nil {
			t.Fatalf("should be fail when password is wrong")
		}
		he = err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategoryUser {
			t.Fatalf("the error category should be user")
		}
	})

	t.Run("login success then logout", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		c := e.NewContext(r, w)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := service.NewUserSession(sess)
		token := util.RandomString(loginTokenLength)
		us.SetLoginToken(token)
		context.SetUserSession(c, us)

		context.SetRequestBody(c, []byte(`{
			"account": "`+account+`",
			"password": "`+util.Sha256(token+password)+`"
		}`))

		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}

		err := ctrl.doLogin(c)
		if err != nil {
			t.Fatalf("login fail, %v", err)
		}

		userInfo := context.GetBody(c).(*UserInfoRes)
		if userInfo.Anonymous || userInfo.Account != account {
			t.Fatalf("login fail")
		}

		// 退出登录
		w = httptest.NewRecorder()
		c = e.NewContext(r, w)
		context.SetUserSession(c, us)
		err = ctrl.doLogout(c)
		if err != nil {
			t.Fatalf("logout fail, %v", err)
		}
		userInfo = context.GetBody(c).(*UserInfoRes)
		if !userInfo.Anonymous || userInfo.Account != "" {
			t.Fatalf("logout fail")
		}
	})

	t.Run("refresh", func(t *testing.T) {
		c := e.NewContext(nil, nil)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := service.NewUserSession(sess)
		context.SetUserSession(c, us)
		err := ctrl.refresh(c)
		if err != nil {
			t.Fatalf("refresh fail, %v", err)
		}
		if us.GetUpdatedAt() == "" {
			t.Fatalf("the updated at should not be empty")
		}
	})
}
