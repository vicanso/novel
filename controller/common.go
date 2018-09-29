package controller

import (
	"github.com/kataras/iris"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
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
	common.Add("GET", "/v1/ip-location", ctrl.getLocationByIP)
}

// getLocationByIP get location by ip
func (c *CommonCtrl) getLocationByIP(ctx iris.Context) {
	query := util.GetRequestQuery(ctx)
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
	setCache(ctx, "1m")
	res(ctx, info)
}
