package controller

import (
	"bytes"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/asset"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xerror"
)

const (
	staticAdminURLPrefix = "/admin/static/"
	staticWebURLPrefix   = "/web/static/"
	staticErrCategory    = "staticServe"
	minCompressSize      = 1024
	defaultIndexFile     = "index.html"
)

var (
	textRegexp = regexp.MustCompile("text|javascript|json")
)

func init() {
	router.Add(
		"GET",
		staticAdminURLPrefix+"*",
		createServe(asset.GetAdminAsset()),
	)
	router.Add(
		"GET",
		staticWebURLPrefix+"*",
		createServe(asset.GetWebAsset()),
	)

	router.Add(
		"GET",
		"/admin/",
		createIndexHandler(asset.GetAdminAsset()),
	)

	router.Add(
		"GET",
		"/web/",
		createIndexHandler(asset.GetWebAsset()),
	)
}

func getNotFoundError(file string) error {
	return &xerror.HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    file + " not found",
		Category:   staticErrCategory,
	}
}

func createServe(as *asset.Asset) echo.HandlerFunc {
	// serve static serve
	return func(c echo.Context) (err error) {
		file := c.Param("*")
		buf := as.Get(file)
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
}

func createIndexHandler(as *asset.Asset) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		file := defaultIndexFile
		buf := as.Get(file)
		context.SetContentType(c, mime.TypeByExtension(filepath.Ext(file)))
		defaultEnv := []byte(`"development"`)
		newEnv := []byte(`"` + config.GetENV() + `"`)
		buf = bytes.Replace(buf, defaultEnv, newEnv, 1)
		context.Res(c, buf)
		return
	}
}
