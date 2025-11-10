package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"goadmin/config"
	"goadmin/internal/context"
	"goadmin/pkg/redisx"
	"goadmin/pkg/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenService struct {
	config *config.JWTConfig
}

// NewTokenService 创建一个新的token服务实例
func NewJwtTokenService(cfg *config.JWTConfig) *JwtTokenService {
	return &JwtTokenService{
		config: cfg,
	}
}

func (s *JwtTokenService) logPrefix() string {
	return "jwt-token-service"
}

// GenerateJWTTokenPair 生成JWT访问令牌和刷新令牌对
func (s *JwtTokenService) GenerateJWTTokenPair(ctx *context.Context, claims Claims) (*TokenPair, error) {
	accessToken, expiresAt, err := s.generateJWT(claims)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshExpire := s.config.RefreshExpire
	if refreshExpire == 0 {
		refreshExpire = 7 * 24 * time.Hour // 默认7天
	}

	refreshToken := util.GenerateUUIDWithoutHyphen()

	// 将刷新令牌存储到Redis中
	refreshKey := s.getRefreshTokenKey(refreshToken)
	// 存储用户ID和访问令牌的相关信息，用于刷新时生成新的访问令牌
	refreshData := claims.String()

	err = redisx.GetClient().Set(ctx, refreshKey, refreshData, refreshExpire).Err()
	if err != nil {
		ctx.Logger.Errorf("%s 构建新的令牌失败: %s %s %v", s.logPrefix(), refreshKey, refreshData, err)
		return nil, fmt.Errorf("保存刷新令牌失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshJWTToken 刷新JWT访问令牌
func (s *JwtTokenService) RefreshJWTToken(
	ctx *context.Context, refreshToken string, f func(Claims) (Claims, error)) (*TokenPair, error) {
	// 验证刷新令牌是否存在
	refreshKey := s.getRefreshTokenKey(refreshToken)
	str, err := redisx.GetClient().Get(ctx, refreshKey).Result()
	if err != nil {
		ctx.Logger.Errorf("%s 刷新令牌无效或已过期: %s %v", s.logPrefix(), refreshKey, err)
		return nil, fmt.Errorf("刷新令牌无效或已过期: %w", err)
	}
	var old Claims
	err = json.Unmarshal([]byte(str), &old)
	if err != nil {
		return nil, fmt.Errorf("解析刷新令牌失败: %w", err)
	}

	newClaims, err := f(old)
	if err != nil {
		ctx.Logger.Errorf("%s 构建新的令牌失败: %s %+v %v", s.logPrefix(), refreshKey, old, err)
		return nil, fmt.Errorf("构建新的令牌失败: %w", err)
	}

	// 生成新的令牌对
	return s.GenerateJWTTokenPair(ctx, newClaims)
}

// ValidateJWTToken 验证JWT令牌并返回claims
func (s *JwtTokenService) ValidateJWTToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// InvalidateRefreshToken 使刷新令牌失效
func (s *JwtTokenService) InvalidateRefreshToken(ctx *context.Context, refreshToken string) error {
	refreshKey := s.getRefreshTokenKey(refreshToken)
	return redisx.GetClient().Del(ctx, refreshKey).Err()
}

// 生成JWT令牌
func (s *JwtTokenService) generateJWT(claims Claims) (string, int64, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 获取过期时间
	expiresAt := claims.ExpiresAt.Unix()

	// 使用密钥签名
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// 获取刷新令牌的Redis键
func (s *JwtTokenService) getRefreshTokenKey(refreshToken string) string {
	prefix := s.config.RefreshTokenKey
	if prefix == "" {
		prefix = "refresh_token:"
	}
	return prefix + refreshToken
}
