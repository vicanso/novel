package middleware

import (
	"testing"
	"time"

	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/session"

	"github.com/labstack/echo"
)

func TestIsLogined(t *testing.T) {
	t.Run("logined", func(t *testing.T) {
		fn := IsLogined(func(c echo.Context) error {
			return nil
		})
		e := echo.New()
		c := e.NewContext(nil, nil)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		err := us.SetAccount("vicanso")
		if err != nil {
			t.Fatalf("set account fail, %v", err)
		}
		err = fn(c)
		if err != nil {
			t.Fatalf("is logined check fail, %v", err)
		}
	})

	t.Run("is not logined", func(t *testing.T) {
		fn := IsLogined(func(c echo.Context) error {
			return nil
		})
		e := echo.New()
		c := e.NewContext(nil, nil)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		err := fn(c)
		if err != errNeedLogined {
			t.Fatalf("should return error when is not logined")
		}
	})
}

func TestIsAnonymous(t *testing.T) {
	t.Run("is not anonymous", func(t *testing.T) {
		fn := IsAnonymous(func(c echo.Context) error {
			return nil
		})
		e := echo.New()
		c := e.NewContext(nil, nil)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		us.SetAccount("vicanso")
		err := fn(c)
		if err != errLoginedAlready {
			t.Fatalf("should return error when is logined")
		}
	})

	t.Run("is anonymous", func(t *testing.T) {
		fn := IsAnonymous(func(c echo.Context) error {
			return nil
		})
		e := echo.New()
		c := e.NewContext(nil, nil)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		us := &service.UserSession{
			Sess: sess,
		}
		context.SetUserSession(c, us)
		err := fn(c)
		if err != nil {
			t.Fatalf("is anonymous check fail, %v", err)
		}
	})
}

func TestWaitFor(t *testing.T) {
	started := time.Now().UnixNano()
	wait := time.Second
	fn := WaitFor(wait)(func(c echo.Context) error {
		return nil
	})
	e := echo.New()
	c := e.NewContext(nil, nil)
	err := fn(c)
	if err != nil {
		t.Fatalf("wait for fail, %v", err)
	}
	use := time.Now().UnixNano() - started
	if use < int64(wait) {
		t.Fatalf("wait for fail")
	}
}
