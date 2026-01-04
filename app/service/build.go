package service

import (
	"fiber-ee/app/service/admin/test"

	"go.uber.org/dig"
)

func BuildAdminServices(c *dig.Container) {
	for _, svc := range adminServices {
		_ = c.Provide(svc)
	}
}

func BuildAppServices(c *dig.Container) {
	for _, svc := range appServices {
		_ = c.Provide(svc)
	}
}

var adminServices = []any{
	test.NewTestService,
}

var appServices = []any{}
