package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
)

func init() {
	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 只针对panic以上的日志增加stack trace
	l, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}
	defaultLogger = l
}

const (
	// LogCategory log category field
	LogCategory = "category"

	// LogAccess access log category
	LogAccess = "access"
	// LogTracker tracker log category
	LogTracker = "tracker"
	// LogUser user log category
	LogUser = "user"
)

// Logger get logger
func Logger() *zap.Logger {
	return defaultLogger
}

// AccessLogger get access logger
func AccessLogger() *zap.Logger {
	return defaultLogger.With(
		zap.String(LogCategory, LogAccess),
	)
}

// TrackerLogger get tracker logger
func TrackerLogger() *zap.Logger {
	return defaultLogger.With(
		zap.String(LogCategory, LogTracker),
	)
}

// UserLogger get user logger
func UserLogger() *zap.Logger {
	return defaultLogger.With(
		zap.String(LogCategory, LogUser),
	)
}
