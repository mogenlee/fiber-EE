package bootstrap

import (
	"context"
	"fmt"

	"fiber-ee/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
	storageRedis "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// NewRedis 创建 Redis 客户端，如果未启用则返回 nil
func NewRedis(cfg *config.Config, log *zap.Logger) *redis.Client {
	if !cfg.Redis.Enabled {
		log.Info("Redis 已禁用，将使用内存存储作为替代")
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Warn("Redis 连接失败 (跳过，可能未重启 Redis 服务)", zap.Error(err))
	} else {
		log.Info("Redis 连接成功")
	}

	return rdb
}

// NewRedisStorage 创建存储后端，如果 rdb 为 nil 则返回内存存储
func NewRedisStorage(rdb *redis.Client) fiber.Storage {
	if rdb == nil {
		return memory.New()
	}
	return storageRedis.NewFromConnection(rdb)
}
