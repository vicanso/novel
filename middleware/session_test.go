package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"

	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/xerror"
)

func TestNewSession(t *testing.T) {
	defaultDuration := time.Hour * 24
	sessConfig := SessionConfig{
		// session cache expires
		Expires: config.GetDurationDefault("session.expires", defaultDuration),
		// the sesion cookie
		Cookie: config.GetSessionCookie(),
		// cookie max age
		CookieMaxAge: config.GetDurationDefault("session.cookie.maxAge", defaultDuration),
		// cookie path
		CookiePath: config.GetCookiePath(),
		// cookie signed keys
		Keys: config.GetSessionKeys(),
	}
	client := service.GetRedisClient()
	t.Run("get session", func(t *testing.T) {
		id := "01CNBNBMNBW92044KPDB8VYKYY"
		buf := []byte(`{
			"a": 1,
			"b": "c"
		}`)
		cmd := client.Set(id, buf, time.Second)
		_, err := cmd.Result()
		if err != nil {
			t.Fatalf("set cache fail, %v", err)
		}
		fn := NewSession(client, sessConfig)(func(c echo.Context) error {
			return nil
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/users/v1/me", nil)
		r.AddCookie(&http.Cookie{
			Name:  "sess",
			Value: id,
		})
		r.AddCookie(&http.Cookie{
			Name:  "sess.sig",
			Value: "rIQ8cMXGRLC22aZeQoU0nZb3BGQ",
		})
		e := echo.New()
		c := e.NewContext(r, w)
		err = fn(c)
		if err != nil {
			t.Fatalf("get session middleware fail, %v", err)
		}
		us := context.GetUserSession(c)
		sess := us.Sess
		if sess.GetInt("a") != 1 || sess.GetString("b") != "c" {
			t.Fatalf("get session data fail")
		}
	})

	t.Run("fetch session fail", func(t *testing.T) {
		id := "01CNBNBMNBW92044KPDB8VYKYY"
		// 非标准json，会导致parse error
		buf := []byte(`{
			"a": 1,
			"b": "c",
		}`)
		cmd := client.Set(id, buf, time.Second)
		_, err := cmd.Result()
		if err != nil {
			t.Fatalf("set cache fail, %v", err)
		}
		fn := NewSession(client, sessConfig)(func(c echo.Context) error {
			return nil
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/users/v1/me", nil)
		r.AddCookie(&http.Cookie{
			Name:  "sess",
			Value: id,
		})
		r.AddCookie(&http.Cookie{
			Name:  "sess.sig",
			Value: "rIQ8cMXGRLC22aZeQoU0nZb3BGQ",
		})
		e := echo.New()
		c := e.NewContext(r, w)
		err = fn(c)
		he := err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategorySession ||
			he.StatusCode != http.StatusInternalServerError {
			t.Fatalf("session fetch error is invalid")
		}
	})

	t.Run("commit session fail", func(t *testing.T) {
		id := "01CNBNBMNBW92044KPDB8VYKYY"
		buf := []byte(`{
			"a": 1,
			"b": "c"
		}`)
		cmd := client.Set(id, buf, time.Second)
		_, err := cmd.Result()
		if err != nil {
			t.Fatalf("set cache fail, %v", err)
		}
		fn := NewSession(client, sessConfig)(func(c echo.Context) error {
			us := context.GetUserSession(c)
			sess := us.Sess
			sess.Set("a", 1)
			return client.Close()
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/users/v1/me", nil)
		r.AddCookie(&http.Cookie{
			Name:  "sess",
			Value: id,
		})
		r.AddCookie(&http.Cookie{
			Name:  "sess.sig",
			Value: "rIQ8cMXGRLC22aZeQoU0nZb3BGQ",
		})
		e := echo.New()
		c := e.NewContext(r, w)
		err = fn(c)
		he := err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategorySession ||
			he.StatusCode != http.StatusInternalServerError {
			t.Fatalf("session commit error is invalid")
		}
	})
}
