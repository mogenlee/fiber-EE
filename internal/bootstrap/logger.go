package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fiber-ee/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	var encoder zapcore.Encoder

	// 定制 Encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置为人可读的 ISO8601 时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		devConfig := zap.NewDevelopmentEncoderConfig()
		devConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(devConfig)
	}

	var level zapcore.Level
	switch cfg.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 默认输出到 Stdout
	writer := zapcore.AddSync(os.Stdout)

	// 如果配置了文件输出且不是 stdout
	if cfg.Log.Output != "" && cfg.Log.Output != "stdout" {
		logDir := cfg.Log.Output
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("无法创建日志目录: %v\n", err)
		} else {
			logFileName := time.Now().Format("2006-01-02") + ".log"
			logPath := filepath.Join(logDir, logFileName)

			file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("无法打开日志文件: %v\n", err)
			} else {
				writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file))
			}
		}
	}

	core := zapcore.NewCore(
		encoder,
		writer,
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger.Info("日志模块初始化成功", zap.String("level", cfg.Log.Level), zap.String("format", cfg.Log.Format))
	return logger
}
