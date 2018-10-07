package middleware

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/xerror"
	"go.uber.org/zap"
)

type (
	// RespondConfig the config of respond middleware
	RespondConfig struct {
	}
)

// NewRespond create a new entry middleware
func NewRespond(config RespondConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := context.GetLogger(c)
			err := next(c)
			// 处理出错的响应
			if err != nil {
				he, ok := err.(*xerror.HTTPError)
				if !ok {
					he = &xerror.HTTPError{
						StatusCode: http.StatusInternalServerError,
						// 如果非http error认为非主动返回异常
						Exception: true,
						Message:   err.Error(),
						Category:  xerror.ErrCategoryCommon,
					}
				}
				status := he.StatusCode
				if status == 0 {
					status = http.StatusInternalServerError
				}
				buf, err := json.Marshal(he)
				if err == nil {
					err = c.JSONBlob(status, buf)
				}
				if err != nil {
					logger.Error("c.JSON fail",
						zap.Error(err),
					)
				}
				return nil
			}
			body := context.GetBody(c)
			status := context.GetStatus(c)
			if body == nil {
				err = c.NoContent(status)
			} else {
				switch body.(type) {
				case string:
					err = c.String(status, body.(string))
				case []byte:
					contentType := context.GetContentType(c)
					if contentType == "" {
						contentType = cs.ContentBinaryHeaderValue
					}
					err = c.Blob(status, contentType, body.([]byte))
				default:
					buf, err := json.Marshal(body)
					if err == nil {
						err = c.JSONBlob(status, buf)
					}
				}
			}
			if err != nil {
				logger.Error("c.JSON fail",
					zap.Error(err),
				)
			}
			return nil
		}
	}
}
