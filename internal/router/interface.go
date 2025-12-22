package router

import "github.com/gofiber/fiber/v3"

// AppRouter 所有功能模块路由必须实现的接口
type AppRouter interface {
	Register(root fiber.Router)
}
