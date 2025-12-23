package middleware

import (
	"time"

	"fiber-ee/config"
	"fiber-ee/internal/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// NewLimiter 创建限流中间件
func NewLimiter(cfg config.Limiter, storage fiber.Storage) fiber.Handler {
	// 未启用则跳过
	if !cfg.Enabled {
		return func(c fiber.Ctx) error {
			return c.Next()
		}
	}

	expiration := time.Duration(cfg.Expiration) * time.Second
	if expiration == 0 {
		expiration = time.Minute
	}

	return limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			if cfg.SkipLocal {
				return c.IP() == "127.0.0.1" || c.IP() == "::1"
			}
			return false
		},
		Max:        cfg.Max,
		Expiration: expiration,
		KeyGenerator: func(c fiber.Ctx) string {
			if xff := c.Get("X-Forwarded-For"); xff != "" {
				return xff
			}
			if xri := c.Get("X-Real-IP"); xri != "" {
				return xri
			}
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return response.TooManyRequests(c)
		},
		Storage: storage,
	})
}
