package router

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
)

func TestAdd(t *testing.T) {
	fn := func(c echo.Context) error {
		return nil
	}
	testPath := "/test-path"
	Add(http.MethodGet, testPath, fn)
	r := routerList[1]
	if r.Method != http.MethodGet ||
		r.Path != testPath ||
		len(r.Mids) != 0 {
		t.Fatalf("add router fail")
	}
}

func TestGroup(t *testing.T) {
	isLogin := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return nil
		}
	}
	g := NewGroup("/users", isLogin)
	getUserOrders := func(c echo.Context) error {
		return nil
	}
	g.Add(http.MethodGet, "/me/orders", getUserOrders)
	r := routerList[2]
	if r.Method != http.MethodGet ||
		r.Path != "/users/me/orders" ||
		len(r.Mids) != 1 {
		t.Fatalf("add group router fail")
	}
}

func TestList(t *testing.T) {
	if len(List()) != len(routerList) {
		t.Fatalf("list function fail")
	}
}
