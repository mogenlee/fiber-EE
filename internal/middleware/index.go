package middleware

import (
	"fiber-ee/internal/pkg/i18n"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/zap"
)

// Use 注册全局中间件（不含认证）
// 包含：panic 恢复、请求ID、请求日志、i18n、安全头、压缩、favicon、跨域
func Use(app *fiber.App, logger *zap.Logger) {
	// panic 恢复，防止单个请求崩溃影响整个服务
	app.Use(recoverErr(logger))

	// 为每个请求生成唯一 ID，便于日志追踪
	app.Use(requestid.New())

	// 请求日志中间件
	app.Use(requestLogger(logger))

	// i18n 国际化中间件
	app.Use(i18n.New())

	// 安全响应头（XSS、点击劫持等防护）
	app.Use(helmet.New())

	// 响应压缩（gzip/brotli）
	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	// favicon 处理，避免 404
	app.Use(favicon.New())

	// 跨域资源共享
	app.Use(cors.New())
}
