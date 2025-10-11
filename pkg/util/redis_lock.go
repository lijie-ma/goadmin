package util

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	client *redis.Client
	key    string
	value  string
	ttl    time.Duration
}

func NewRedisLock(client *redis.Client, key, value string, ttl time.Duration) *RedisLock {
	return &RedisLock{
		client: client,
		key:    key,
		value:  value,
		ttl:    ttl,
	}
}

// TryLock 尝试获取锁
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
	ok, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
	return ok, err
}

// Unlock 释放锁（必须比对 value，避免误删别人的锁）
func (l *RedisLock) Unlock(ctx context.Context) (bool, error) {
	// 使用 Lua 脚本保证原子性：先比对 value，再删除
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)
	res, err := script.Run(ctx, l.client, []string{l.key}, l.value).Result()
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}
