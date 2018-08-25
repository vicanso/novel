package controller

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/utils"
)

type (
	// LocationByIPParams params of location by ip
	LocationByIPParams struct {
		IP   string `valid:"ipv4"`
		Type string `valid:"int,in(1|2)"`
	}
	// commonCtrl common controller
	commonCtrl struct {
	}
)

func init() {
	ctrl := commonCtrl{}
	common := router.NewGroup("/common")
	common.Add("GET", "/v1/ip-location", ctrl.getLocationByIP)
}

// getLocationByIP get location by ip
func (c *commonCtrl) getLocationByIP(ctx iris.Context) {
	query := utils.GetRequestQuery(ctx)
	params := &LocationByIPParams{}
	err := validate(params, query)
	if err != nil {
		resErr(ctx, err)
		return
	}
	info, err := service.GetLocationByIP(params.IP)
	if err != nil {
		resErr(ctx, err)
		return
	}
	setCache(ctx, "10m")
	res(ctx, info)
}
