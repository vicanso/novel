package controller

import (
	"github.com/labstack/echo"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/validate"
)

type (
	// LocationByIPParams params of location by ip
	LocationByIPParams struct {
		IP string `valid:"ipv4"`
	}
	// CommonCtrl common controller
	CommonCtrl struct {
	}
)

func init() {
	ctrl := CommonCtrl{}
	common := router.NewGroup("/common")

	common.Add(
		"GET",
		"/v1/ip-location",
		ctrl.getLocationByIP,
	)
}

// getLocationByIP get location by ip
func (cc *CommonCtrl) getLocationByIP(c echo.Context) (err error) {
	query := getRequestQuery(c)
	params := LocationByIPParams{}
	err = validate.Do(&params, query)
	if err != nil {
		return
	}
	info, err := service.GetLocationByIP(params.IP)
	if err != nil {
		return
	}
	setCache(c, "1m")
	res(c, info)
	return
}
