package main

import (
	"context"
	"fiber-ee/internal/bootstrap"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// @title Fiber-EE API
// @version 1.0.0
// @description 企业级 Go Web 框架 API 文档
// @host localhost:8080
// @BasePath /admin/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	c := bootstrap.BuildContainer()

	// 注册路由
	_ = c.Provide(bootstrap.NewServer)

	err := c.Invoke(func(s *bootstrap.Server) {
		// 启动服务器
		go func() {
			if err := s.Start(); err != nil {
				s.Log.Fatal("服务器运行失败", zap.Error(err))
			}
		}()

		// 打印已注册路由
		time.Sleep(100 * time.Millisecond)
		for _, r := range s.App.GetRoutes(true) {
			if r.Method == fiber.MethodHead {
				continue
			}
			fmt.Printf("[%-4s] %s %s\n", r.Method, r.Path, r.Name)
		}

		// 等待中断信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		sig := <-quit

		s.Log.Info("收到关闭信号", zap.String("signal", sig.String()))

		// 优雅关闭
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			if err := s.Shutdown(); err != nil {
				s.Log.Error("服务器关闭出错", zap.Error(err))
			}
			close(done)
		}()

		select {
		case <-done:
			s.Log.Info("服务器已优雅关闭")
		case <-ctx.Done():
			s.Log.Warn("服务器关闭超时，强制退出")
		}
	})

	if err != nil {
		panic(err)
	}
}
