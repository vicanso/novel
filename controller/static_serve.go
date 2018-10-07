package controller

import (
	"mime"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/asset"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xerror"
)

const (
	staticURLPrefix   = "/static/"
	staticErrCategory = "staticServe"
	minCompressSize   = 1024
)

var (
	textRegexp = regexp.MustCompile("text|javascript|json")
)

func init() {
	router.Add(
		"GET",
		staticURLPrefix+"*",
		serve,
	)
}

func getNotFoundError(file string) error {
	return &xerror.HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    file + " not found",
		Category:   staticErrCategory,
	}
}

// serve static serve
func serve(c echo.Context) (err error) {
	file := c.Param("*")
	buf := asset.Get(file)
	bufSize := len(buf)
	if bufSize == 0 {
		err = getNotFoundError(file)
		return
	}
	contentType := mime.TypeByExtension(filepath.Ext(file))
	// 如果是文本类，而且数据大于最小压缩长度，则执行压缩
	if textRegexp.MatchString(contentType) && bufSize > minCompressSize {
		buf, err = util.Gzip(buf, 0)
		if err != nil {
			return
		}
		context.SetHeader(c, echo.HeaderContentEncoding, cs.Gzip)
	}
	context.SetContentType(c, contentType)
	context.Res(c, buf)
	return
}
