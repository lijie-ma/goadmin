package user

import (
	"goadmin/internal/context"
	"goadmin/internal/repository/user"
	"goadmin/internal/service/token"
	userSrv "goadmin/internal/service/user"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户相关的API路由
func RegisterRoutes(r *gin.RouterGroup) {

	// 创建依赖服务
	userRepo := user.NewUserRepository()
	tokenSrv := token.NewTokenService()
	userService := userSrv.NewUserService()
	handler := NewHandler(userService, userRepo, tokenSrv)

	group := r.Group("/user")
	{
		// 登录接口 - 不需要认证
		group.POST("/login", context.Build(handler.Login))

		// 需要认证的接口
		authGroup := group.Group("/")
		// TODO: 添加认证中间件
		// authGroup.Use(middleware.JWT())
		{
			authGroup.GET("/logout", context.Build(handler.Logout))
			authGroup.POST("/changePwd", context.Build(handler.ChangePassword))
		}
	}
}
