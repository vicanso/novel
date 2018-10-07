package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xerror"
)

var (
	errNeedLogined = &xerror.HTTPError{
		StatusCode: http.StatusUnauthorized,
		Category:   xerror.ErrCategoryUser,
		Message:    "please login first",
	}
	errLoginedAlready = &xerror.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   xerror.ErrCategoryUser,
		Message:    "user is logined, please logout first",
	}
	errFunctionForbidden = &xerror.HTTPError{
		StatusCode: http.StatusForbidden,
		Category:   xerror.ErrCategoryUser,
		Message:    "the function is forbidden",
	}
)

// IsLogined check login status，if not, will return error
func IsLogined(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		us := context.GetUserSession(c)
		if us == nil || us.GetAccount() == "" {
			return errNeedLogined
		}
		return next(c)
	}
}

// IsAnonymous check login status, if yes, will return error
func IsAnonymous(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		us := context.GetUserSession(c)
		if us == nil || us.GetAccount() != "" {
			return errLoginedAlready
		}
		return next(c)
	}
}

// IsSu check the user roles include su
func IsSu(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		us := context.GetUserSession(c)
		if us == nil || us.GetAccount() == "" {
			return errNeedLogined
		}
		roles := us.GetRoles()
		if !util.ContainsString(roles, "su") {
			return errFunctionForbidden
		}
		return next(c)
	}
}

// WaitFor at least wait for duration
func WaitFor(d time.Duration) echo.MiddlewareFunc {
	ns := d.Nanoseconds()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			start := time.Now()
			err = next(c)
			use := time.Now().UnixNano() - start.UnixNano()
			if use < ns {
				time.Sleep(time.Duration(ns-use) * time.Nanosecond)
			}
			// 无论成功还是失败都wait for
			return
		}
	}
}
