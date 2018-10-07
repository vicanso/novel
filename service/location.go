package service

import (
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/helper"
	"github.com/vicanso/novel/xerror"
)

var (
	// ErrIPLocationNotFound ip location not found error
	ErrIPLocationNotFound = xerror.New("IP Location not found")
)

type (
	// IPLocation ip location
	IPLocation struct {
		IP      string `json:"ip"`
		Country string `json:"country"`
		Region  string `json:"region"`
		ISP     string `json:"isp"`
	}
)

// GetLocationByIP get location by ip
func GetLocationByIP(ip string) (info *IPLocation, err error) {
	url := config.GetString("locationByIP")
	buf, err := helper.HTTPGet(url, map[string]string{
		"ip": ip,
	})
	if err != nil {
		return
	}
	if len(buf) == 0 {
		err = ErrIPLocationNotFound
		return
	}
	str := json.Get(buf, "data").ToString()
	info = &IPLocation{}
	err = json.UnmarshalFromString(str, info)
	return
}
