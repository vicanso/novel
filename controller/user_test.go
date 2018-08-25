package controller

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kataras/iris"

	"github.com/vicanso/novel/utils"
	"github.com/vicanso/session"
)

func TestUserCtrl(t *testing.T) {
	ctrl := userCtrl{}
	cookies := []string{}
	t.Run("getInfo", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/users/v1/me", nil)
		w := httptest.NewRecorder()
		sess := &session.Session{}
		ctx := utils.NewContext(w, r)
		utils.SetSession(ctx, sess)
		ctrl.getInfo(ctx)

		cookies = ctx.ResponseWriter().Header()["Set-Cookie"]

		userInfo := utils.GetBody(ctx).(*userInfoResponse)
		if !userInfo.Anonymous {
			t.Fatalf("user info should be anonymous")
		}

		if userInfo.Date == "" {
			t.Fatalf("user info's date should not be empty")
		}
	})

	t.Run("getAvatar", func(t *testing.T) {
		ctx := utils.NewResContext()
		ctrl.getAvatar(ctx)
		if !strings.HasPrefix(ctx.GetContentType(), "image/jpeg") {
			t.Fatalf("the content type should be jpeg")
		}
		buf := utils.GetBody(ctx).([]byte)

		if base64.StdEncoding.EncodeToString(buf) != avatar {
			t.Fatalf("the content is wrong")
		}
	})

	t.Run("getLoginToken", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/users/v1/me/token", nil)
		w := httptest.NewRecorder()
		sess := &session.Session{}
		ctx := utils.NewContext(w, r)
		utils.SetSession(ctx, sess)
		ctrl.getLoginToken(ctx)
		data := utils.GetBody(ctx).(iris.Map)
		if len(data["token"].(string)) != 8 {
			t.Fatalf("get login token fail")
		}
	})

	t.Run("doLogin/Logout", func(t *testing.T) {
		params := `{
			"account": "vicanso",
			"password": "treexie"
		}`
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me", nil)
		w := httptest.NewRecorder()

		// 参数校验出错
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		ctx := utils.NewContext(w, r)
		utils.SetRequestBody(ctx, []byte("{}"))
		ctrl.doLogin(ctx)
		errData := utils.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest || errData["category"] != utils.ErrCategoryValidate {
			t.Fatalf("login params should be invalid")
		}

		// 无track cookie的出错
		utils.SetRequestBody(ctx, []byte(params))
		utils.SetSession(ctx, sess)
		ctrl.doLogin(ctx)
		errData = utils.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest ||
			errData["message"] != "track key is not found" {
			t.Fatalf("no track key should be error")
		}

		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}
		ctrl.doLogin(ctx)
		data := utils.GetBody(ctx).(*userInfoResponse)
		if data.Account != "vicanso" {
			t.Fatalf("login fail")
		}

		// login already
		ctrl.doLogin(ctx)
		errData = utils.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest ||
			errData["message"] != "account is logined, please logout first" {
			t.Fatalf("login already should be error")
		}

		ctrl.doLogout(ctx)
		data = utils.GetBody(ctx).(*userInfoResponse)
		if !data.Anonymous {
			t.Fatalf("logout fail")
		}
	})

	t.Run("refresh", func(t *testing.T) {
		ctx := utils.NewResContext()
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		utils.SetSession(ctx, sess)
		ctrl.refresh(ctx)
		if ctx.GetStatusCode() != http.StatusNoContent {
			t.Fatalf("http status should be 204")
		}
		if sess.GetUpdatedAt() == "" {
			t.Fatalf("the updated at should not be empty")
		}

	})
}
