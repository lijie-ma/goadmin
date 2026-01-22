package middleware

import (
	"goadmin/config"
	"goadmin/pkg/logger"
	"goadmin/pkg/trace"
	"net/http"
	"strings"

	"goadmin/internal/i18n"
	modeluser "goadmin/internal/model/user"
	"goadmin/internal/repository/role"
	userrepo "goadmin/internal/repository/user"
	tokenService "goadmin/internal/service/token"

	"github.com/gin-gonic/gin"
)

var (
	// 白名单
	whiteList = []string{}

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
			abortWithError(
				c,
				http.StatusUnauthorized,
				i18n.E(c, "common.NotFound", map[string]any{"item": i18n.T(c, "common.item.token", nil)}))
			return
		}

		// 解析Token头部
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && (parts[0] == "Bearer" || parts[0] == "bearer")) {
			abortWithError(
				c,
				http.StatusUnauthorized,
				i18n.E(c, "common.NotFound", map[string]any{"item": i18n.T(c, "common.item.token", nil)}))
			return
		}

		tokenSrv := tokenService.NewJwtTokenService(&config.Get().JWT)
		claims, err := tokenSrv.ValidateJWTToken(parts[1])
		if err != nil || !claims.IsAdmin() {
			logger.Global().With(trace.GetTrace(c)).Errorf("token invalid %v", err)
			abortWithError(
				c,
				http.StatusUnauthorized,
				i18n.E(c, "common.InvalidParameter", map[string]any{"item": i18n.T(c, "common.item.token", nil)}))
			return
		}
		sessionData, err := generateUserSession(c, claims.UserID)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, err)
			return
		}
		err = hasPermission(c, sessionData)
		if err != nil {
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
		logger.Global().With(trace.GetTrace(ctx)).Errorf("generateUserSession GetByID %v", err)
		return nil, i18n.E(ctx, "common.RepositoryErr", nil)
	}
	if !u.IsActive() {
		return nil, i18n.E(ctx, "user.AccountStatusAbnormal", nil)
	}
	return u, nil
}

// 是否有权限
func hasPermission(ctx *gin.Context, u *modeluser.User) error {
	if u.IsSuperAdmin() {
		return nil
	}
	path := strings.TrimLeft(ctx.Request.URL.Path, "/")
	hasPerm, err := role.NewRolePermissionRepositoryWithDB().HasAccessURL(ctx, u.RoleCode, path)
	if err != nil {
		logger.Global().With(trace.GetTrace(ctx)).Errorf("hasPermission HasAccessURL %v", err)
		return i18n.E(ctx, "common.RepositoryErr", nil)
	}
	if !hasPerm {
		return i18n.E(ctx, "common.PermissionDeny", nil)
	}
	return nil
}

// abortWithError 中止请求并返回错误
func abortWithError(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": err.Error(),
	})
}
