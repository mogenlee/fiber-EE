package bootstrap

import (
	"fiber-ee/internal/service/admin/test"

	"go.uber.org/dig"
)

func buildAdminServices(c *dig.Container) {
	for _, svc := range adminServices {
		_ = c.Provide(svc)
	}
}

func buildAppServices(c *dig.Container) {
	for _, svc := range appServices {
		_ = c.Provide(svc)
	}
}

var adminServices = []any{
	test.NewTestService,
}

var appServices = []any{}
