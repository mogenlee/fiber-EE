package middleware

import (
	"errors"
	"fiber-ee/internal/pkg/response"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// NewErrorHandler 创建全局错误处理器
// 功能：
//   - 统一处理所有未捕获的错误
//   - 自动识别自定义业务错误（RespType）和系统错误
//   - 记录错误日志（业务错误 warn 级别，系统错误 error 级别）
//   - 返回统一格式的 JSON 响应
func NewErrorHandler(log *zap.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}

		// 获取请求开始时间（用于计算耗时）
		var latency time.Duration
		if start, ok := c.Locals("request_start").(time.Time); ok {
			latency = time.Since(start)
		}

		// 从 response header 获取 request_id
		requestID := c.GetRespHeader("X-Request-ID")

		// 识别自定义业务错误类型
		var resp response.RespType
		if errors.As(err, &resp) {
			// 业务错误，warn 级别日志
			log.Warn("业务错误",
				zap.Int("code", resp.Code()),
				zap.String("msg", resp.Msg()),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.Duration("latency", latency),
				zap.String("ip", c.IP()),
				zap.String("request_id", requestID),
			)
			return response.Fail(c, resp)
		}

		// 系统错误，error 级别日志
		log.Error("系统错误",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Duration("latency", latency),
			zap.String("ip", c.IP()),
			zap.String("request_id", requestID),
		)
		return response.Fail(c, response.ErrInternal)
	}
}
