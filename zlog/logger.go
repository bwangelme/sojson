package zlog

import (
	"context"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger 初始化日志系统
func InitLogger(logLevel string, logPath string) error {
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 如果指定了日志路径，则写入文件，否则输出到控制台
	if logPath != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// 打开日志文件
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Logger.SetOutput(logFile)
	} else {
		Logger.SetOutput(os.Stdout)
	}

	return nil
}

// Debug 记录 Debug 级别日志
func Debug(ctx context.Context, args ...interface{}) {
	Logger.WithContext(ctx).Debug(args...)
}

// Debugf 记录格式化的 Debug 级别日志
func Debugf(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Debugf(format, args...)
}

// Info 记录 Info 级别日志
func Info(ctx context.Context, args ...interface{}) {
	Logger.WithContext(ctx).Info(args...)
}

// Infof 记录格式化的 Info 级别日志
func Infof(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Infof(format, args...)
}

// Warn 记录 Warn 级别日志
func Warn(ctx context.Context, args ...interface{}) {
	Logger.WithContext(ctx).Warn(args...)
}

// Warnf 记录格式化的 Warn 级别日志
func Warnf(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Warnf(format, args...)
}

// Error 记录 Error 级别日志
func Error(ctx context.Context, args ...interface{}) {
	Logger.WithContext(ctx).Error(args...)
}

// Errorf 记录格式化的 Error 级别日志
func Errorf(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Errorf(format, args...)
}

// Fatal 记录 Fatal 级别日志并退出程序
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf 记录格式化的 Fatal 级别日志并退出程序
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// WithField 添加字段到日志
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

// WithFields 添加多个字段到日志
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}
