package bootstrap

import (
	"context"
	"fiber-ee/internal/model/entity"
	"fiber-ee/internal/model/query"
	"fmt"
	"time"

	"fiber-ee/config"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDatabase(cfg *config.Config, log *zap.Logger) (*gorm.DB, *query.Query) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.DBName,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)
		dialector = postgres.Open(dsn)
	default:
		// 默认使用 sqlite
		dialector = sqlite.Open(cfg.Database.Source)
	}

	// GORM 配置
	gormConfig := &gorm.Config{}

	// 表前缀配置
	if cfg.Database.TablePrefix != "" {
		gormConfig.NamingStrategy = schema.NamingStrategy{
			TablePrefix: cfg.Database.TablePrefix,
		}
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Fatal("数据库连接失败", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("获取 DB 实例失败", zap.Error(err))
	} else {
		sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
	}

	// 自动迁移
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Error("数据库自动迁移失败", zap.Error(err))
	}

	// 初始化 gorm-gen query（必须在使用 query.User 之前）
	query.SetDefault(db)

	log.Info("数据库初始化成功",
		zap.String("driver", cfg.Database.Driver),
		zap.String("prefix", cfg.Database.TablePrefix),
	)

	// 插入测试数据（仅开发环境）
	if cfg.App.Debug {
		q := query.User.WithContext(context.Background())
		count, _ := q.Count()
		if count == 0 {
			// 批量生成 21 条测试数据
			users := make([]*entity.User, 21)
			for i := 0; i < 21; i++ {
				users[i] = &entity.User{
					Username: fmt.Sprintf("test%d", i+1),
					Password: "123456",
					Email:    fmt.Sprintf("test%d@qq.com", i+1),
				}
			}

			if err := q.Create(users...); err != nil {
				log.Error("插入测试数据失败", zap.Error(err))
			} else {
				log.Info("已插入测试用户", zap.Int("count", len(users)))
			}
		}
	}

	return db, query.Q
}
