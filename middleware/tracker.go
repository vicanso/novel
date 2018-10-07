package middleware

import (
	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
)

const (
	// HandleSuccess handle success
	HandleSuccess = iota
	// HandleFail handle fail
	HandleFail
)

type (
	// OnTrack on track function
	OnTrack func(*TrackerInfo, echo.Context)
	// TrackerConfig tracker config
	TrackerConfig struct {
		OnTrack OnTrack
	}
	// TrackerInfo tracker info
	TrackerInfo struct {
		Query  map[string]string
		Params map[string]string
		Form   []byte
		Result int
		Err    error
	}
)

// NewTracker create a tracker middleware
func NewTracker(config TrackerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			info := &TrackerInfo{
				Query:  context.GetRequestQuery(c),
				Params: context.GetRequestParams(c),
				Form:   context.GetRequestBody(c),
			}
			err = next(c)
			if err != nil {
				info.Result = HandleFail
				info.Err = err
			}

			if config.OnTrack != nil {
				config.OnTrack(info, c)
			}
			return
		}
	}
}
