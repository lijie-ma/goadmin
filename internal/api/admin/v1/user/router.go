package user

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"
	"goadmin/internal/repository/user"
	"goadmin/internal/service/token"
	userSrv "goadmin/internal/service/user"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户相关的API路由
func RegisterRoutes(
	r *gin.RouterGroup,
	userService userSrv.UserService,
	userRepository user.UserRepository,
	tokenService *token.TokenService,
) {
	handler := NewHandler(userService, userRepository, tokenService)

	group := r.Group("/user")
	{
		// 登录接口 - 不需要认证
		group.POST("/login", context.Build(handler.Login))

		// 需要认证的接口
		authGroup := group.Group("/")

		authGroup.Use(middleware.Auth())
		{
			authGroup.GET("/logout", context.Build(handler.Logout))
			authGroup.GET("/info", context.Build(handler.GetCurrentUser))
			authGroup.POST("/change_pwd", context.Build(handler.ChangePassword))
			authGroup.POST("/reset_pwd", context.Build(handler.ResetPassword))
			authGroup.GET("/list", context.Build(handler.ListUsers))
			authGroup.POST("/create", context.Build(handler.CreateUser))
			authGroup.POST("/update", context.Build(handler.UpdateUser))
			authGroup.POST("/delete", context.Build(handler.DeleteUser))
		}
	}
}