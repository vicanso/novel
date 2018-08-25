package service

import (
	"net/http"

	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/utils"
)

var (
	// ErrIPLocationNotFound ip location not found error
	ErrIPLocationNotFound = &utils.HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "IP Location not found",
	}
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
	buf, err := utils.HTTPGet(url, map[string]string{
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
