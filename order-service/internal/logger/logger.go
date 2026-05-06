package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init 初始化日志
func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 设置时间格式
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	return nil
}

// Info 记录信息日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Error 记录错误日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Warn 记录警告日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Debug 记录调试日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Sync 刷新日志缓冲
func Sync() {
	_ = Log.Sync()
}
