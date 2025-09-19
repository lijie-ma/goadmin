package redis

import (
	"context"
	"goadmin/config"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	// 加载配置
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}
	t.Logf("配置加载成功")
	t.Logf("Redis地址: %s:%d", cfg.Redis.Host, cfg.Redis.Port)

	// 初始化Redis
	t.Log("开始初始化Redis...")
	err = Init(cfg)
	if err != nil {
		t.Fatalf("初始化Redis失败: %v", err)
	}
	t.Log("Redis初始化成功")

	// 验证Redis连接
	client := GetClient()
	if client == nil {
		t.Fatal("Redis客户端实例为空")
	}
	t.Log("Redis客户端实例获取成功")

	// 测试Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Log("尝试Ping Redis...")
	result, err := client.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Redis Ping失败: %v", err)
	}
	t.Logf("Redis Ping成功: %s", result)

	// 尝试简单的SET/GET操作
	testKey := "test:redis:conn"
	testValue := "hello-redis"

	t.Logf("尝试设置键: %s", testKey)
	err = client.Set(ctx, testKey, testValue, 1*time.Minute).Err()
	if err != nil {
		t.Fatalf("Redis SET操作失败: %v", err)
	}

	t.Logf("尝试获取键: %s", testKey)
	val, err := client.Get(ctx, testKey).Result()
	if err != nil {
		t.Fatalf("Redis GET操作失败: %v", err)
	}

	if val != testValue {
		t.Fatalf("Redis值不匹配，期望: %s，实际: %s", testValue, val)
	}

	t.Logf("Redis SET/GET操作成功")
	t.Log("Redis连接测试成功")
}
