package middleware

import (
	"errors"
	"fmt"
	"goadmin/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTConfig JWT中间件配置
type JWTConfig struct {
	// 密钥
	SecretKey string
	// Token在请求头中的键名
	TokenHeadName string
	// 从Token中解析出的用户ID在上下文中的键名
	IdentityKey string
	// 忽略认证的路径
	ExcludePaths []string
}

// DefaultJWTConfig 默认JWT配置
func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:     "goadmin_secret_key",
		TokenHeadName: "Authorization",
		IdentityKey:   "user_id",
		ExcludePaths:  []string{"/api/v1/auth/login", "/api/v1/auth/register", "/health"},
	}
}

// JWT 返回一个JWT认证中间件
func JWT(config ...JWTConfig) gin.HandlerFunc {
	// 使用默认配置或传入的配置
	var cfg JWTConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultJWTConfig()
	}

	return func(c *gin.Context) {
		// 检查是否为忽略认证的路径
		path := c.Request.URL.Path
		for _, excludePath := range cfg.ExcludePaths {
			if strings.HasPrefix(path, excludePath) {
				c.Next()
				return
			}
		}

		// 从请求头获取Token
		authHeader := c.GetHeader(cfg.TokenHeadName)
		if authHeader == "" {
			abortWithError(c, http.StatusUnauthorized, errors.New("未提供认证信息"))
			return
		}

		// 解析Token头部
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && (parts[0] == "Bearer" || parts[0] == "bearer")) {
			abortWithError(c, http.StatusUnauthorized, errors.New("认证格式无效"))
			return
		}

		token := parts[1]

		// 验证Token并解析声明
		// 注意：这里只是一个示例，实际项目中应该实现完整的JWT验证逻辑
		claims, err := parseToken(token, cfg.SecretKey)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, err)
			return
		}

		// 将用户ID存入上下文
		c.Set(cfg.IdentityKey, claims.UserID)

		// 继续处理请求
		c.Next()
	}
}

// TokenClaims JWT声明结构
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	ExpireAt int64  `json:"expire_at"`
}

// 解析JWT Token
func parseToken(tokenString, secretKey string) (*TokenClaims, error) {
	// 注意：这只是一个示例框架，实际项目中需要使用jwt-go等库来实现完整的JWT解析逻辑
	// 例如：
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//     return []byte(secretKey), nil
	// })

	// 示例实现（非实际解析逻辑）
	if tokenString == "" {
		return nil, errors.New("空认证Token")
	}

	// 检查Token有效期（示例）
	claims := &TokenClaims{
		UserID:   "demo_user_id", // 在真实实现中，这些值应该从JWT解析得到
		Username: "demo_user",
		Role:     "admin",
		ExpireAt: time.Now().Add(time.Hour).Unix(),
	}

	if claims.ExpireAt < time.Now().Unix() {
		return nil, errors.New("认证Token已过期")
	}

	return claims, nil
}

// abortWithError 中止请求并返回错误
func abortWithError(c *gin.Context, code int, err error) {
	log := logger.Global()
	log.Error("认证失败",
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
		logger.String("error", err.Error()),
	)

	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": fmt.Sprintf("认证失败: %v", err),
	})
}
