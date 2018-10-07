package context

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/service"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xlog"
	"go.uber.org/zap"
)

const (
	// Status 记录status code
	Status = "_status"
	// Body 记录body
	Body = "_body"
	// UserSession 记录user session
	UserSession = "_userSession"
	// RequestBody 记录请求数据
	RequestBody = "_requestBody"
	// RequestQuery 设置请求的query
	RequestQuery = "_requestQuery"
	// Logger 记录track的logger
	Logger = "_logger"
	// RequestID the request
	RequestID = "_requestID"
	// ContentType the content type of response
	ContentType = "_contentType"
)

// Res 设置响应数据
func Res(c echo.Context, data interface{}) {
	if c.Get(Body) != nil {
		// TODO 此处应该增加告警
		logger := GetLogger(c)
		if logger != nil {
			logger.Error("duplicate set body")
		}
	}
	c.Set(Body, data)
}

// GetBody get the body of response
func GetBody(c echo.Context) interface{} {
	return c.Get(Body)
}

// SetStatus set the status of response
func SetStatus(c echo.Context, stauts int) {
	c.Set(Status, stauts)
}

// GetStatus get the status of response
func GetStatus(c echo.Context) (status int) {
	v := c.Get(Status)
	if v != nil {
		status, _ = v.(int)
	} else {
		// 若未设置status的处理
		if c.Get(Body) == nil {
			status = http.StatusNoContent
		} else {
			status = http.StatusOK
		}
	}
	return
}

// SetContentType set the content-type of response
func SetContentType(c echo.Context, t string) {
	c.Set(ContentType, t)
}

// GetContentType get the content-type of response
func GetContentType(c echo.Context) string {
	v := c.Get(ContentType)
	if v == nil {
		return ""
	}
	return v.(string)
}

// ResCreated set the body of response and status to 201
func ResCreated(c echo.Context, data interface{}) {
	SetStatus(c, http.StatusCreated)
	Res(c, data)
}

// SetUserSession set the session
func SetUserSession(c echo.Context, userSession *service.UserSession) {
	c.Set(UserSession, userSession)
}

// GetUserSession get the session
func GetUserSession(c echo.Context) (userSession *service.UserSession) {
	v := c.Get(UserSession)
	if v == nil {
		return
	}
	userSession, _ = v.(*service.UserSession)
	return
}

// SetRequestBody set the request body
func SetRequestBody(c echo.Context, buf []byte) {
	c.Set(RequestBody, buf)
}

// GetRequestBody get the request body
func GetRequestBody(c echo.Context) (buf []byte) {
	v := c.Get(RequestBody)
	if v == nil {
		return
	}
	buf = v.([]byte)
	return
}

// GetRequestQuery get the request query (use map[string]string)
func GetRequestQuery(c echo.Context) map[string]string {
	// 如果已有缓存，则直接返回
	v := c.Get(RequestQuery)
	if v != nil {
		return v.(map[string]string)
	}
	values := c.QueryParams()
	if len(values) == 0 {
		return nil
	}

	m := make(map[string]string)
	for k, v := range values {
		m[k] = v[0]
	}
	c.Set(RequestQuery, m)
	return m
}

// GetRequestParams get request params
func GetRequestParams(c echo.Context) (m map[string]string) {
	params := c.ParamNames()
	if len(params) == 0 {
		return nil
	}
	values := c.ParamValues()
	m = make(map[string]string)
	for index, name := range params {
		m[name] = values[index]
	}
	return
}

// GetTrackID get track id
func GetTrackID(c echo.Context) string {
	if c.Request() == nil {
		return ""
	}
	cookie, _ := c.Cookie(config.GetTrackKey())
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

// GetRequestID get the request id
func GetRequestID(c echo.Context) (id string) {
	v := c.Get(RequestID)
	if v == nil {
		id = util.GenUlid()
		c.Set(RequestID, id)
	} else {
		id = v.(string)
	}
	return
}

// SetHeader set the header of response
func SetHeader(c echo.Context, key, value string) {
	c.Response().Header().Set(key, value)
}

// SetNoCache set the response to be no cache
func SetNoCache(c echo.Context) {
	SetHeader(c, cs.HeaderCacheControl, cs.HeaderNoCache)
}

// SetCache set the max age of response
func SetCache(c echo.Context, age string) error {
	d, err := time.ParseDuration(age)
	if err != nil {
		return err
	}
	cache := "public, max-age=" + strconv.Itoa(int(d.Seconds()))
	SetHeader(c, cs.HeaderCacheControl, cache)
	return nil
}

// SetPrivateCache set the private max age of response
func SetPrivateCache(c echo.Context, age string) error {
	d, err := time.ParseDuration(age)
	if err != nil {
		return err
	}
	cache := "private, max-age=" + strconv.Itoa(int(d.Seconds()))
	SetHeader(c, cs.HeaderCacheControl, cache)
	return nil
}

// SetCacheWithSMaxAge set the cache with s-maxage
func SetCacheWithSMaxAge(c echo.Context, age, sMaxAge string) error {
	dMaxAge, err := time.ParseDuration(age)
	if err != nil {
		return err
	}
	dSMaxAge, err := time.ParseDuration(sMaxAge)
	if err != nil {
		return err
	}
	cache := fmt.Sprintf("public, max-age=%d, s-maxage=%d", int(dMaxAge.Seconds()), int(dSMaxAge.Seconds()))
	SetHeader(c, cs.HeaderCacheControl, cache)
	return nil
}

// GetLogger get the logger from context
func GetLogger(c echo.Context) (logger *zap.Logger) {
	v := c.Get(Logger)
	if v != nil {
		logger = v.(*zap.Logger)
	} else {
		logger = xlog.Logger().With(
			zap.String("trackId", GetTrackID(c)),
			zap.String("requestId", GetRequestID(c)),
		)
		c.Set(Logger, logger)
	}
	return
}
