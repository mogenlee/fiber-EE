package bootstrap

import (
	"fiber-ee/app/router"
	"fiber-ee/app/service"
	"fiber-ee/internal/middleware"
	"fiber-ee/internal/pkg/request"
	"fiber-ee/internal/pkg/validator"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	// 注册配置
	_ = c.Provide(NewConfig)

	// 注册日志
	_ = c.Provide(NewLogger)

	// 注册数据库
	_ = c.Provide(NewDatabase)

	// 注册 Casbin
	_ = c.Provide(NewCasbinEnforcer)

	// 注册 Redis
	_ = c.Provide(NewRedis)
	_ = c.Provide(NewRedisStorage)

	// 注册验证器
	_ = c.Provide(validator.NewValidator)

	// 注册请求绑定器
	_ = c.Provide(request.NewValidate)

	// 注册 JWT 配置
	_ = c.Provide(middleware.NewJWTConfig)

	// 注册 Admin Service
	service.BuildAdminServices(c)
	// 注册 App Service
	service.BuildAppServices(c)

	// 注册 Admin Routers
	router.BuildAdminRoutes(c)

	// 注册 App Routers
	router.BuildAppRoutes(c)

	return c
}
