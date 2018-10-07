package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/xerror"

	"github.com/labstack/echo"

	"github.com/h2non/gock"
)

func TestCommonCtrl(t *testing.T) {
	ctrl := CommonCtrl{}
	t.Run("getLocationByIP", func(t *testing.T) {
		defer gock.Off()
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/ip-location?ip=abcd", nil)
		w := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(r, w)
		err := ctrl.getLocationByIP(c)
		he := err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategoryValidte ||
			he.StatusCode != http.StatusBadRequest {
			t.Fatalf("the error is invalid")
		}

		gock.New("http://ip.taobao.com").
			Get("/service/getIpInfo.php").
			Reply(500).
			BodyString("{}")
		r = httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/ip-location?ip=114.114.114.114", nil)
		w = httptest.NewRecorder()
		c = e.NewContext(r, w)
		err = ctrl.getLocationByIP(c)
		he = err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategoryRequset ||
			he.StatusCode != http.StatusInternalServerError {
			t.Fatalf("the error is invalid")
		}

		gock.New("http://ip.taobao.com").
			Get("/service/getIpInfo.php").
			Reply(200).
			BodyString(`{"code":0,"data":{"ip":"114.114.114.114","country":"中国","area":"","region":"江苏","city":"南京","county":"XX","isp":"XX","country_id":"CN","area_id":"","region_id":"320000","city_id":"320100","county_id":"xx","isp_id":"xx"}}`)

		r = httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/ip-location?ip=114.114.114.114", nil)
		w = httptest.NewRecorder()
		c = e.NewContext(r, w)

		err = ctrl.getLocationByIP(c)
		if err != nil {
			t.Fatalf("get location by ip fail, %v", err)
		}

		buf, err := json.Marshal(context.GetBody(c))
		if err != nil {
			t.Fatalf("response data is not json, %v", err)
		}

		if string(buf) != `{"ip":"114.114.114.114","country":"中国","region":"江苏","isp":"XX"}` {
			t.Fatalf("respons data is wrong")
		}
	})
}
