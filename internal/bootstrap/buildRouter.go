package bootstrap

import (
	"fiber-ee/internal/router"
	"fiber-ee/internal/router/admin/auth"
	"fiber-ee/internal/router/admin/test"

	"go.uber.org/dig"
)

func buildAdminRoutes(c *dig.Container) {
	for _, routers := range adminRouters {
		_ = c.Provide(routers, dig.As(new(router.AppRouter)), dig.Group("admin_routers"))
	}
}

func buildAppRoutes(c *dig.Container) {
	for _, routers := range appRouters {
		_ = c.Provide(routers, dig.As(new(router.AppRouter)), dig.Group("app_routers"))
	}
}

var adminRouters = []any{
	auth.NewRouter,
	test.NewRouter,
}

var appRouters = []any{}
