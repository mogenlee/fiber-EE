package bootstrap

import (
	"fiber-ee/config"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewCasbinEnforcer 初始化 Casbin 执行器
func NewCasbinEnforcer(db *gorm.DB, cfg *config.Config, log *zap.Logger) (*casbin.Enforcer, error) {
	// 使用 GORM 适配器，策略存储在数据库
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 从配置文件加载模型
	enforcer, err := casbin.NewEnforcer(cfg.Casbin.ModelPath, adapter)
	if err != nil {
		return nil, err
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	// 添加默认策略（示例）
	initDefaultPolicies(enforcer)

	log.Info("Casbin 初始化成功", zap.String("model", cfg.Casbin.ModelPath))
	return enforcer, nil
}

// initDefaultPolicies 初始化默认策略
func initDefaultPolicies(e *casbin.Enforcer) {
	// 角色权限：admin 可以访问所有 /admin/v1/* 接口
	if has, _ := e.HasPolicy("admin", "/admin/v1/*", "GET"); !has {
		_, _ = e.AddPolicy("admin", "/admin/v1/*", "GET")
		_, _ = e.AddPolicy("admin", "/admin/v1/*", "POST")
		_, _ = e.AddPolicy("admin", "/admin/v1/*", "PUT")
		_, _ = e.AddPolicy("admin", "/admin/v1/*", "DELETE")
	}

	// 用户角色绑定示例：user_1 是 admin 角色
	if has, _ := e.HasGroupingPolicy("user_1", "admin"); !has {
		_, _ = e.AddGroupingPolicy("user_1", "admin")
	}
}
