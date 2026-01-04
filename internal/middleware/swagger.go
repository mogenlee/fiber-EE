package middleware

import (
	"fiber-ee/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/wdcbot/qingfeng"
)

func swagger(cfg *config.Config) fiber.Handler {
	return adaptor.HTTPHandler(qingfeng.HTTPHandler(qingfeng.Config{
		Title:         cfg.App.Name + " API",
		Description:   "企业级 Go Web 框架 API 文档",
		Version:       "1.0.0",
		BasePath:      "/doc",
		DocPath:       "./docs/swagger.json",
		EnableDebug:   cfg.App.Debug,
		DarkMode:      false,
		UITheme:       qingfeng.ThemeModern,
		AutoGenerate:  true,
		SwagSearchDir: "./",
		SwagOutputDir: "./docs",
		SwagArgs:      []string{"-g", "cmd/server/main.go", "--parseDependency"},
	}))
}
