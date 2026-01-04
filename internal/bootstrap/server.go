package bootstrap

import (
	"context"
	"fiber-ee/app/router"
	"fiber-ee/config"
	"fiber-ee/internal/middleware"
	"fiber-ee/internal/pkg/validator"
	"time"

	"github.com/bytedance/sonic"
	"github.com/casbin/casbin/v2"
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
	Config       *config.Config
	Log          *zap.Logger
	Validator    *validator.CustomValidator
	Redis        *redis.Client
	DB           *gorm.DB
	Storage      fiber.Storage
	Enforcer     *casbin.Enforcer
	AdminRouters []router.AppRouter `group:"admin_routers"`
	AppRouters   []router.AppRouter `group:"app_routers"`
}

// NewServer 创建并初始化服务器实例
func NewServer(p NewServerParams) *Server {
	app := fiber.New(fiber.Config{
		AppName:         p.Config.App.Name,
		JSONDecoder:     sonic.Unmarshal,
		JSONEncoder:     sonic.Marshal,
		StructValidator: p.Validator,
		ErrorHandler:    middleware.NewErrorHandler(p.Log),
	})
	// 注册监控路由
	app.Get("/metrics", monitor.New(monitor.Config{Title: p.Config.App.Name}))

	// 注册全局中间件
	middleware.Use(app, p.Config, p.Log, p.Storage)

	// 注册 Admin 路由组（JWT + Casbin）
	admin := app.Group("/admin/v1") //middleware.JWTAuth(p.Config),
	//middleware.CasbinAuth(p.Enforcer),

	// 注册 Admin 路由
	for _, r := range p.AdminRouters {
		r.Register(admin)
	}

	// 注册 App 路由组（公开接口）
	api := app.Group("/api/v1")
	for _, r := range p.AppRouters {
		r.Register(api)
	}

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
func (s *Server) Shutdown() error {
	if err := s.App.ShutdownWithTimeout(10 * time.Second); err != nil {
		s.Log.Error("HTTP 服务关闭失败", zap.Error(err))
		return err
	}
	s.Log.Info("HTTP 服务已关闭")

	if s.Redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Redis.Close(); err != nil {
			s.Log.Error("Redis 关闭失败", zap.Error(err))
		} else {
			s.Log.Info("Redis 连接已关闭")
		}
		_ = ctx
	}

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
