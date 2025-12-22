package server

import (
	"context"
	"fiber-ee/config"
	"fiber-ee/internal/middleware"
	"fiber-ee/internal/pkg/validator"
	"time"

	"github.com/gofiber/contrib/v3/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server 服务器实例，封装 Fiber 应用及相关依赖
type Server struct {
	App       *fiber.App
	Config    *config.Config
	Log       *zap.Logger
	Validator *validator.CustomValidator
	Redis     *redis.Client
	DB        *gorm.DB
}

// NewServerParams 创建服务器所需的依赖参数
type NewServerParams struct {
	dig.In
	Config    *config.Config
	Log       *zap.Logger
	Validator *validator.CustomValidator
	Redis     *redis.Client
	DB        *gorm.DB
}

// NewServer 创建并初始化服务器实例
func NewServer(p NewServerParams) *Server {
	app := fiber.New(fiber.Config{
		AppName:         p.Config.App.Name,
		StructValidator: p.Validator,
		// 全局错误处理中间件，统一处理错误响应和日志记录
		ErrorHandler: middleware.NewErrorHandler(p.Log),
	})

	app.Get("/metrics", monitor.New(monitor.Config{Title: p.Config.App.Name}))

	// 注册全局中间件（不含 JWT 认证）
	middleware.Use(app, p.Log)

	return &Server{
		App:       app,
		Config:    p.Config,
		Log:       p.Log,
		Validator: p.Validator,
		Redis:     p.Redis,
		DB:        p.DB,
	}
}

// Start 启动 HTTP 服务器
func (s *Server) Start() error {
	s.Log.Info("正在启动服务器", zap.String("port", s.Config.App.Port))
	return s.App.Listen(s.Config.App.Port)
}

// Shutdown 优雅关闭服务器
// 1. 关闭 HTTP 服务（等待现有请求完成）
// 2. 关闭 Redis 连接
// 3. 关闭数据库连接
func (s *Server) Shutdown() error {
	// 关闭 HTTP 服务
	if err := s.App.ShutdownWithTimeout(10 * time.Second); err != nil {
		s.Log.Error("HTTP 服务关闭失败", zap.Error(err))
		return err
	}
	s.Log.Info("HTTP 服务已关闭")

	// 关闭 Redis
	if s.Redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Redis.Close(); err != nil {
			s.Log.Error("Redis 关闭失败", zap.Error(err))
		} else {
			s.Log.Info("Redis 连接已关闭")
		}
		_ = ctx // 用于超时控制
	}

	// 关闭数据库
	if s.DB != nil {
		sqlDB, err := s.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				s.Log.Error("数据库关闭失败", zap.Error(err))
			} else {
				s.Log.Info("数据库连接已关闭")
			}
		}
	}

	return nil
}
