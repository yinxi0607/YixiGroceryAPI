package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// Logger 是全局日志实例
var Logger *logrus.Logger

// Init 初始化日志配置
func Init(logFilePath string, logLevel string) error {
	Logger = logrus.New()

	// 设置日志格式为 JSON
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// 设置日志级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // 默认级别
	}
	Logger.SetLevel(level)

	// 配置日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // 天
		Compress:   true,
	}

	// 设置日志输出（文件和控制台）
	Logger.SetOutput(io.MultiWriter(lumberjackLogger, os.Stdout))

	return nil
}
