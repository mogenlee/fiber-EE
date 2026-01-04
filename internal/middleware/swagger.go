package middleware

import (
	"fiber-ee/config"
	"net/http"

	"github.com/wdcbot/qingfeng"
)

func swagger(cfg *config.Config) http.Handler {
	return qingfeng.HTTPHandler(qingfeng.Config{
		Title:         cfg.App.Name + " API",
		Description:   "企业级 Go Web 框架 API 文档",
		Version:       "1.0.0",
		BasePath:      "/doc",
		DocPath:       "./internal/docs/swagger.json",
		EnableDebug:   cfg.App.Debug,
		DarkMode:      false,
		UITheme:       qingfeng.ThemeModern,
		AutoGenerate:  true,
		SwagSearchDir: "./",
		SwagOutputDir: "./internal/docs",
		SwagArgs:      []string{"-g", "cmd/server/main.go", "--parseDependency"},
	})
}
