package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/session"
)

const (
	// Body 记录body
	Body = "_body"
	// Session 记录session
	Session = "_session"
	// RequestBody 记录请求数据
	RequestBody = "_requestBody"
	// RequestQuery 设置请求的query
	RequestQuery = "_requestQuery"
	// Logger 记录track的logger
	Logger = "_logger"

	// AccountField session account field
	AccountField = "account"

	// HeaderCacheControl cache-control header
	HeaderCacheControl = "Cache-Control"
)

var (
	stackReg = regexp.MustCompile(`\((0x[\s\S]+)\)`)
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Res 响应数据
func Res(ctx iris.Context, data interface{}) {
	ctx.Values().Set(Body, data)
}

// ResNoContent 返回无内容(204)
func ResNoContent(ctx iris.Context) {
	ctx.StatusCode(http.StatusNoContent)
	Res(ctx, nil)
}

// ResCreated 响应为201
func ResCreated(ctx iris.Context, data interface{}) {
	ctx.StatusCode(http.StatusCreated)
	Res(ctx, data)
}

// ResJPEG 响应jpeg
func ResJPEG(ctx iris.Context, data []byte) {
	ctx.ContentType("image/jpeg")
	Res(ctx, data)
}

// ResPNG 响应png
func ResPNG(ctx iris.Context, data []byte) {
	ctx.ContentType("image/png")
	Res(ctx, data)
}

// ResErr 出错处理
func ResErr(ctx iris.Context, err error) {
	data := iris.Map{
		"message": err.Error(),
	}
	status := http.StatusInternalServerError
	// 根据error类型生成各类的状态码与出错信息
	if he, ok := err.(*HTTPError); ok {
		status = he.StatusCode
		if status == 0 {
			status = http.StatusInternalServerError
		}
		if he.Code != "" {
			data["code"] = he.Code
		}
		if he.Category != "" {
			data["category"] = he.Category
		}
		if he.Extra != nil {
			data["extra"] = he.Extra
		}
	}
	if !IsProduction() {
		data["stack"] = GetStack(2 << 10)
	}
	SetNoCache(ctx)
	ctx.StatusCode(status)
	Res(ctx, data)
}

// GetBody 获取响应数据
func GetBody(ctx iris.Context) interface{} {
	return ctx.Values().Get(Body)
}

// SetSession 设置session
func SetSession(ctx iris.Context, sess *session.Session) {
	ctx.Values().Set(Session, sess)
}

// GetSession 获取session
func GetSession(ctx iris.Context) (sess *session.Session) {
	v := ctx.Values().Get(Session)
	if v == nil {
		return
	}
	sess, _ = v.(*session.Session)
	return
}

// GetAccount get account info
func GetAccount(ctx iris.Context) string {
	s := GetSession(ctx)
	if s == nil {
		return ""
	}
	return s.GetString(AccountField)
}

// GetTrackID get track id
func GetTrackID(ctx iris.Context) string {
	return ctx.GetCookie(config.GetTrackKey())
}

// SetRequestBody 设置请求数据
func SetRequestBody(ctx iris.Context, buf []byte) {
	ctx.Values().Set(RequestBody, buf)
	return
}

// GetRequestBody 获取请求数据
func GetRequestBody(ctx iris.Context) (buf []byte) {
	v := ctx.Values().Get(RequestBody)
	if v == nil {
		return
	}
	buf = v.([]byte)
	return
}

// GetRequestQuery get query string
// 默认生成的query为map[string][]string，因为标准的query允许相同的参数，会生成数组，但实际使用中不常使用相同参数，为了方便开发，因此增加了parse为map[string]string的处理
func GetRequestQuery(ctx iris.Context) map[string]string {
	url := ctx.Request().URL
	if url.RawQuery == "" {
		return nil
	}
	v := ctx.Values().Get(RequestQuery)
	if v != nil {
		return v.(map[string]string)
	}
	m := make(map[string]string)
	q := url.Query()
	for k, v := range q {
		m[k] = v[0]
	}
	ctx.Values().Set(RequestQuery, m)
	return m
}

// Validate 校验数据
func Validate(s interface{}, data interface{}) (err error) {
	// TODO 了解自定义的校验方式
	if data != nil {
		switch data.(type) {
		case []byte:
			err = json.Unmarshal(data.([]byte), s)
			if err != nil {
				fmt.Println(err)
				err = NewJSONParseError(err)
				return
			}
		default:
			buf, e := json.Marshal(data)
			if e != nil {
				err = NewJSONParseError(e)
				return
			}
			e = json.Unmarshal(buf, s)
			if e != nil {
				err = NewJSONParseError(e)
				return
			}
		}
	}
	_, err = govalidator.ValidateStruct(s)
	if err != nil {
		err = NewValidateError(err)
	}
	return
}

// GetStack 获取调用信息
func GetStack(size int) []string {
	stack := make([]byte, size)
	runtime.Stack(stack, true)
	arr := strings.Split(string(stack), "\n")
	max := len(arr)
	result := []string{}
	// goroutine与此函数的stack无需展示，因此index从3开始
	for index := 3; index < max; index += 2 {
		if index+1 >= max {
			break
		}
		tmpArr := strings.Split(arr[index], "/")
		fn := stackReg.ReplaceAllString(tmpArr[len(tmpArr)-1], "")
		str := fn + ": " + strings.Replace(arr[index+1], "\t", "", 1)
		result = append(result, str)
	}
	return result
}

// SetHeader 设置Header（覆盖非添加）
func SetHeader(ctx iris.Context, key, value string) {
	header := ctx.ResponseWriter().Header()
	header.Set(key, value)
}

// SetNoCache 设置无缓存
func SetNoCache(ctx iris.Context) {
	SetHeader(ctx, HeaderCacheControl, "no-cache, max-age=0")
}

// SetNoStore 设置不可保存
func SetNoStore(ctx iris.Context) {
	SetHeader(ctx, HeaderCacheControl, "no-store")
}

// SetCache 设置缓存
func SetCache(ctx iris.Context, age string) error {
	d, err := time.ParseDuration(age)
	if err != nil {
		return err
	}
	cache := "public, max-age=" + strconv.Itoa(int(d.Seconds()))
	SetHeader(ctx, HeaderCacheControl, cache)
	return nil
}

// NewContext 创建一个新的context
func NewContext(w http.ResponseWriter, r *http.Request) iris.Context {
	app := iris.New()
	return app.ContextPool.Acquire(w, r)
}

// NewResContext 创建带response的context
func NewResContext() iris.Context {
	return NewContext(httptest.NewRecorder(), nil)
}
