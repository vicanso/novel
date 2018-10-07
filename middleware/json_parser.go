package middleware

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/xerror"
)

const (
	// 默认为50kb
	defaultRequestJSONLimit = 50 * 1024
)

type (
	// JSONParserConfig json parser middleware config
	JSONParserConfig struct {
		Limit int
	}
)

var (
	// errJSONTooLarge too large error
	errJSONTooLarge = &xerror.HTTPError{
		StatusCode: http.StatusRequestEntityTooLarge,
		Message:    "request post json too large",
		Category:   xerror.ErrCategoryIO,
	}
)

// NewJSONParser create a new json parser middleware
func NewJSONParser(config JSONParserConfig) echo.MiddlewareFunc {
	limit := defaultRequestJSONLimit
	if config.Limit != 0 {
		limit = config.Limit
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			method := req.Method
			if method != http.MethodPost && method != http.MethodPatch && method != http.MethodPut {
				return next(c)
			}
			contentType := req.Header.Get("Content-Type")
			// 非json则跳过
			if !strings.HasPrefix(contentType, "application/json") {
				return next(c)
			}
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				err = xerror.NewIO(err.Error())
				return
			}
			if limit != 0 && len(body) > limit {
				err = errJSONTooLarge
				return
			}
			context.SetRequestBody(c, body)
			return next(c)
		}
	}
}
