package bootstrap

import (
	"fiber-ee/internal/handler/admin/auth"
	"fiber-ee/internal/handler/admin/test"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/dig"
)

// AppRouter 路由接口
type AppRouter interface {
	Register(root fiber.Router)
}

func buildAdminRoutes(c *dig.Container) {
	for _, routers := range adminRouters {
		_ = c.Provide(routers, dig.As(new(AppRouter)), dig.Group("admin_routers"))
	}
}

func buildAppRoutes(c *dig.Container) {
	for _, routers := range appRouters {
		_ = c.Provide(routers, dig.As(new(AppRouter)), dig.Group("app_routers"))
	}
}

var adminRouters = []any{
	auth.NewRouter,
	test.NewRouter,
}

var appRouters = []any{}
