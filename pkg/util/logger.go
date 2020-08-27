package util

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CustomLogger struct
type CustomLogger struct {
	prefix string
	logger *zap.Logger
}

// NewCustomLogger func
func NewCustomLogger(prefix, level string) *CustomLogger {
	logger := GetLogger(level, "console", false, 1)
	return &CustomLogger{prefix: prefix, logger: logger}
}

func (l *CustomLogger) Error(msg string) {
	l.logger.Sugar().Error(l.withPrefix(msg))
}

// Infof func
func (l *CustomLogger) Infof(msg string, args ...interface{}) {
	l.logger.Sugar().Infof(l.withPrefix(msg), args...)
}

// Debugf func
func (l *CustomLogger) Debugf(msg string, args ...interface{}) {
	l.logger.Sugar().Debugf(l.withPrefix(msg), args...)
}

func (l *CustomLogger) withPrefix(msg string) string {
	return fmt.Sprintf("%s: %s", l.prefix, msg)
}

// GetLogger returns a logger
func GetLogger(level, encoding string, showFullCaller bool, callerSkip int) *zap.Logger {
	zapCfg := zap.NewDevelopmentConfig()

	zapCfg.EncoderConfig.TimeKey = "time"
	zapCfg.EncoderConfig.LevelKey = "level"
	zapCfg.EncoderConfig.NameKey = "logger"
	zapCfg.EncoderConfig.CallerKey = "caller"
	zapCfg.EncoderConfig.MessageKey = "msg"
	zapCfg.EncoderConfig.StacktraceKey = "stacktrace"
	zapCfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if showFullCaller {
		zapCfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}
	// update user change
	// update log level
	switch strings.ToLower(level) {
	case "debug":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warning":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	// update encoding type
	switch strings.ToLower(encoding) {
	case "json":
		zapCfg.Encoding = "json"
	default:
		zapCfg.Encoding = "console"
	}

	logger, err := zapCfg.Build(zap.AddCaller(), zap.AddCallerSkip(callerSkip))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return logger
}
