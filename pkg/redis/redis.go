package redis

import (
	"context"
	"fmt"
	"goadmin/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis 客户端实例
var Client *redis.Client

// 初始化Redis客户端
func Init(cfg *config.RedisConfig) error {
	// 如果Redis未启用，直接返回
	if !cfg.Enable {
		return nil
	}

	// 创建Redis客户端
	Client = NewClient(cfg, cfg.DB)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("Redis连接测试失败: %w", err)
	}

	return nil
}

// GetClient 获取Redis客户端实例
func GetClient() *redis.Client {
	return Client
}

// NewClient 创建Redis客户端实例
func NewClient(cfg *config.RedisConfig, dbSelect int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              dbSelect,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.IdleTimeout,
		ConnMaxLifetime: cfg.MaxConnAge,
	})
}
