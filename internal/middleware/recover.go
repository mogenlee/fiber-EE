package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
)

func recoverErr(logger *zap.Logger) fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			logger.Error("panic recovered",
				zap.String("panic", fmt.Sprintf("%v", e)),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("ip", c.IP()),
			)
		},
	})
}
