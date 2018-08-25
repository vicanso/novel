package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/session"

	"github.com/kataras/iris"

	"github.com/vicanso/novel/utils"
)

func TestIsLogined(t *testing.T) {
	t.Run("not logined", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := utils.NewContext(w, r)
		IsLogined(ctx)
		if ctx.GetStatusCode() != http.StatusUnauthorized {
			t.Fatalf("the status code should be 401")
		}
		err := utils.GetBody(ctx).(iris.Map)
		if err["category"].(string) != utils.ErrCategoryLogic ||
			err["code"].(string) != utils.ErrCodeUser {
			t.Fatalf("the http error is not wrong")
		}
	})

	t.Run("logined", func(t *testing.T) {
		sess := session.Mock(session.M{
			"fetched": true,
			"data": session.M{
				"account": "vicanso",
			},
		})
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := utils.NewContext(w, r)
		utils.SetSession(ctx, sess)
		IsLogined(ctx)
		fmt.Println(ctx.GetStatusCode())
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("is login check fail")
		}
	})
}
