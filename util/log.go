package util

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/kataras/iris"
	"go.uber.org/zap"
)

var (
	defaultLogger *zap.Logger
	stackReg      = regexp.MustCompile(`\((0x[\s\S]+)\)`)
	appPath       = ""
)

func init() {
	c := zap.NewProductionConfig()
	if IsDevelopment() {
		c = zap.NewDevelopmentConfig()
	}
	l, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}
	defaultLogger = l
	_, file, _, _ := runtime.Caller(0)
	fileDivide := "/"
	arr := strings.Split(file, fileDivide)
	arr = arr[0 : len(arr)-2]
	appPath = strings.Join(arr, fileDivide) + fileDivide
}

const (
	// LogCategory log category field
	LogCategory = "category"
	// LogTrack log track field
	LogTrack = "track"

	// LogAccess access log category
	LogAccess = "access"
	// LogTracker tracker log category
	LogTracker = "tracker"
	// LogUser user log category
	LogUser = "user"
)

// GetLogger get logger
func GetLogger() *zap.Logger {
	return defaultLogger
}

// CreateAccessLogger 创建access logger
func CreateAccessLogger() *zap.Logger {
	return defaultLogger.With(zap.String(LogCategory, LogAccess))
}

// CreateTrackerLogger 创建tracker logger
func CreateTrackerLogger() *zap.Logger {
	return defaultLogger.With(zap.String(LogCategory, LogTracker))
}

// CreateUserLogger 创建user logger
func CreateUserLogger(ctx iris.Context) *zap.Logger {
	return defaultLogger.With(
		zap.String(LogCategory, LogUser),
		zap.String(LogTrack, GetTrackID(ctx)))
}

// SetContextLogger 设置logger
func SetContextLogger(ctx iris.Context, logger *zap.Logger) {
	ctx.Values().Set(Logger, logger)
}

// GetContextLogger 获取logger
func GetContextLogger(ctx iris.Context) *zap.Logger {
	logger := ctx.Values().Get(Logger)
	if logger == nil {
		return nil
	}
	return logger.(*zap.Logger)
}

// GetStack 获取调用信息
func GetStack(start, end int) []string {
	size := 2 << 10
	stack := make([]byte, size)
	runtime.Stack(stack, true)
	arr := strings.Split(string(stack), "\n")
	arr = arr[1:]
	max := len(arr) - 1
	result := []string{}
	for index := 0; index < max; index += 2 {
		if index+1 >= max {
			break
		}

		file := strings.Replace(arr[index+1], appPath, "", 1)
		tmpArr := strings.Split(arr[index], "/")
		fn := stackReg.ReplaceAllString(tmpArr[len(tmpArr)-1], "")
		str := fn + ": " + file
		result = append(result, str)
	}
	return result[start:end]
}
