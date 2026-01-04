package middleware

import (
	"fiber-ee/config"
	"fiber-ee/internal/pkg/i18n"

	"github.com/gofiber/contrib/v3/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/zap"
)

// Use 注册全局中间件（不含认证）
func Use(app *fiber.App, cfg *config.Config, logger *zap.Logger, storage fiber.Storage) {
	//	swagger
	app.Use(swagger(cfg))

	// 监控路由
	app.Use("/metrics", monitor.New(monitor.Config{Title: cfg.App.Name}))

	// panic 恢复
	app.Use(recoverErr(logger))

	// 请求 ID
	app.Use(requestid.New())

	// 请求日志
	app.Use(requestLogger(logger))

	// i18n 国际化
	app.Use(i18n.New())

	// 安全响应头
	app.Use(helmet.New())

	// 响应压缩
	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	// favicon
	app.Use(favicon.New())

	// 跨域
	app.Use(cors.New())

	// 限流
	app.Use(NewLimiter(cfg.Limiter, storage))

}
