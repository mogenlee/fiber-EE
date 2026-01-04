package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// requestLogger 请求日志中间件
// 只记录成功请求，错误请求由 ErrorHandler 统一记录
func requestLogger(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		// 存储开始时间，供 ErrorHandler 计算耗时
		c.Locals("request_start", start)

		// 执行后续处理
		err := c.Next()

		// 有错误时跳过日志（由 ErrorHandler 统一记录）
		if err != nil {
			return err
		}

		// 记录成功请求日志
		duration := time.Since(start)
		// 从 response header 获取 request_id（requestid 中间件设置）
		requestID := c.GetRespHeader("X-Request-ID")
		logger.Info("请求完成",
			zap.String("method", c.Method()),
			zap.String("name", c.Route().Name),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", duration),
			zap.String("ip", c.IP()),
			zap.String("request_id", requestID),
		)

		return nil
	}
}
