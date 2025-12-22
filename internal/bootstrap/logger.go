package bootstrap

import (
	"os"
	"path/filepath"
	"time"

	"fiber-ee/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// dateWriter 按日期切割的日志写入器
type dateWriter struct {
	logDir     string
	currentDay string
	file       *lumberjack.Logger
}

func newDateWriter(logDir string) *dateWriter {
	w := &dateWriter{logDir: logDir}
	w.rotate()
	return w
}

func (w *dateWriter) Write(p []byte) (n int, err error) {
	today := time.Now().Format("2006-01-02")
	if today != w.currentDay {
		w.rotate()
	}
	return w.file.Write(p)
}

func (w *dateWriter) rotate() {
	w.currentDay = time.Now().Format("2006-01-02")
	w.file = &lumberjack.Logger{
		Filename:   filepath.Join(w.logDir, w.currentDay+".log"),
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     30, // days
		Compress:   true,
	}
}

func (w *dateWriter) Sync() error {
	return nil
}

func NewLogger(cfg *config.Config) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		devConfig := zap.NewDevelopmentEncoderConfig()
		devConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(devConfig)
	}

	level := parseLogLevel(cfg.Log.Level)

	// 默认输出到 Stdout
	writer := zapcore.AddSync(os.Stdout)

	// 配置文件输出
	if cfg.Log.Output != "" && cfg.Log.Output != "stdout" {
		_ = os.MkdirAll(cfg.Log.Output, 0755)
		dateWriter := newDateWriter(cfg.Log.Output)
		writer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(dateWriter),
		)
	}

	core := zapcore.NewCore(encoder, writer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger.Info("日志模块初始化成功", zap.String("level", cfg.Log.Level), zap.String("format", cfg.Log.Format))
	return logger
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
