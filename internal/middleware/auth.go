package middleware

import (
	"errors"
	"fmt"
	"goadmin/config"
	"goadmin/pkg/logger"
	"goadmin/pkg/trace"
	"net/http"
	"strings"

	modeluser "goadmin/internal/model/user"
	userrepo "goadmin/internal/repository/user"
	tokenService "goadmin/internal/service/token"

	"github.com/gin-gonic/gin"
)

var (
	// 白名单
	whiteList = []string{
		"/api/v1/user/login",
		"/health",
	}
	tokenHeadName = "Authorization"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为忽略认证的路径
		path := c.Request.URL.Path
		for _, excludePath := range whiteList {
			if strings.HasPrefix(path, excludePath) {
				c.Next()
				return
			}
		}

		// 从请求头获取Token
		authHeader := c.GetHeader(tokenHeadName)
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

		tokenSrv := tokenService.NewJwtTokenService(&config.Get().JWT)
		claims, err := tokenSrv.ValidateJWTToken(parts[1])
		if err != nil || !claims.IsAdmin() {
			abortWithError(c, http.StatusUnauthorized, err)
			return
		}
		sessionData, err := generateUserSession(c, claims.UserID)
		if err != nil || !claims.IsAdmin() {
			abortWithError(c, http.StatusUnauthorized, err)
			return
		}
		// 将用户存入上下文
		c.Set(gin.AuthUserKey, sessionData)

		// 继续处理请求
		c.Next()
	}
}

func generateUserSession(ctx *gin.Context, userID uint64) (*modeluser.User, error) {
	u, err := userrepo.NewUserRepository().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !u.IsActive() {
		return nil, fmt.Errorf("账户状态异常")
	}
	return u, nil
}

// abortWithError 中止请求并返回错误
func abortWithError(c *gin.Context, code int, err error) {
	log := logger.Global().With(trace.GetTrace(c))
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
