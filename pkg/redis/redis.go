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

	// 解析超时配置
	dialTimeout, err := time.ParseDuration(cfg.DialTimeout.String())
	if err != nil {
		return fmt.Errorf("解析Redis.DialTimeout失败: %w", err)
	}

	readTimeout, err := time.ParseDuration(cfg.ReadTimeout.String())
	if err != nil {
		return fmt.Errorf("解析Redis.ReadTimeout失败: %w", err)
	}

	writeTimeout, err := time.ParseDuration(cfg.WriteTimeout.String())
	if err != nil {
		return fmt.Errorf("解析Redis.WriteTimeout失败: %w", err)
	}

	poolTimeout, err := time.ParseDuration(cfg.PoolTimeout.String())
	if err != nil {
		return fmt.Errorf("解析Redis.PoolTimeout失败: %w", err)
	}

	idleTimeout, err := time.ParseDuration(cfg.IdleTimeout.String())
	if err != nil {
		return fmt.Errorf("解析Redis.IdleTimeout失败: %w", err)
	}

	maxConnAge, err := time.ParseDuration(cfg.MaxConnAge.String())
	if err != nil {
		return fmt.Errorf("解析Redis.MaxConnAge失败: %w", err)
	}

	// 创建Redis客户端
	Client = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolTimeout:     poolTimeout,
		ConnMaxIdleTime: idleTimeout,
		ConnMaxLifetime: maxConnAge,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接测试失败: %w", err)
	}

	return nil
}

// GetClient 获取Redis客户端实例
func GetClient() *redis.Client {
	return Client
}
