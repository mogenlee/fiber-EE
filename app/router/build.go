package router

import (
	"fiber-ee/app/router/admin/auth"
	"fiber-ee/app/router/admin/test"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/dig"
)

// AppRouter 路由接口
type AppRouter interface {
	Register(root fiber.Router)
}

func BuildAdminRoutes(c *dig.Container) {
	for _, routers := range adminRouters {
		_ = c.Provide(routers, dig.As(new(AppRouter)), dig.Group("admin_routers"))
	}
}

func BuildAppRoutes(c *dig.Container) {
	for _, routers := range appRouters {
		_ = c.Provide(routers, dig.As(new(AppRouter)), dig.Group("app_routers"))
	}
}

var adminRouters = []any{
	test.NewRouter,
	auth.NewRouter,
}

var appRouters = []any{}
