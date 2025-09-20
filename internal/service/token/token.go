package token

import (
	"goadmin/internal/context"
	"goadmin/pkg/redis"
	"goadmin/pkg/util"
	"time"
)

// TokenService token服务结构体
type TokenService struct {
}

// NewTokenService 创建一个新的token服务实例
func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) logPrefix() string {
	return "token-service"
}

// GenerateToken 生成一个新的token
// expiration为可选参数，如果不传则使用默认过期时间
func (s *TokenService) GenerateToken(ctx *context.Context, expiration ...time.Duration) (string, error) {
	// 生成UUID作为token
	token := util.GenerateUUID()

	// 设置过期时间
	exp := time.Minute // 默认1分钟
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	// 将token存储到Redis中
	key := "token:" + token
	err := redis.GetClient().Set(ctx, key, true, exp).Err()
	if err != nil {
		ctx.Logger.Errorf("%s 构建新的令牌失败: %s %+v", s.logPrefix(), key, err)
		return "", err
	}

	return token, nil
}

// DeleteToken 删除指定的token
func (s *TokenService) DeleteToken(ctx *context.Context, token string) error {
	key := "token:" + token
	return redis.GetClient().Del(ctx, key).Err()
}

// ValidateToken 验证token是否有效
func (s *TokenService) ValidateToken(ctx *context.Context, token string) bool {
	key := "token:" + token
	exists, err := redis.GetClient().Exists(ctx, key).Result()
	return err == nil && exists == 1
}

// ExtendToken 延长token的有效期
func (s *TokenService) ExtendToken(ctx *context.Context, token string, duration time.Duration) error {
	key := "token:" + token
	return redis.GetClient().Expire(ctx, key, duration).Err()
}
