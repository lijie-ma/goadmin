package api

import (
	"goadmin/internal/api/admin/v1/captcha"
	"goadmin/internal/api/admin/v1/operate_log"
	"goadmin/internal/api/admin/v1/position"
	"goadmin/internal/api/admin/v1/role"
	"goadmin/internal/api/admin/v1/setting"
	"goadmin/internal/api/admin/v1/upload"
	userapi "goadmin/internal/api/admin/v1/user"
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/middleware"
	"goadmin/internal/service/token"
	userservice "goadmin/internal/service/user"
	roleservice "goadmin/internal/service/role"
	positionservice "goadmin/internal/service/position"
	operatelogsService "goadmin/internal/service/operate_log"
	settingsservice "goadmin/internal/service/setting"
	"goadmin/internal/repository/user"
	"os"

	"github.com/gin-gonic/gin"
)

// Services holds all services for dependency injection into routers
type Services struct {
	TokenService      *token.TokenService
	UserService       userservice.UserService
	RoleService       roleservice.RoleService
	PositionService   positionservice.PositionService
	OperateLogService operatelogsService.OperateLogService
	SettingService    settingsservice.ServerSettingService
	UserRepository    user.UserRepository
}

func RegisterRouter(r *gin.Engine, services Services) {
	r.Use(
		middleware.Trace(),
		i18n.Middleware(),
		middleware.Logger(),
		middleware.Header(),
		middleware.Recovery(),
	)
	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API路由组
	adminHandler(r, services)
}

func adminHandler(r *gin.Engine, services Services) {
	adminGroup := r.Group("/admin/v1")
	{
		// 验证码路由
		captchaRouter := adminGroup.Group("/captcha")
		{
			captchaRouter.GET("/generate", context.Build(captcha.GenerateHandler))
			captchaRouter.POST("/check", context.Build(captcha.CheckHandler))
		}

		// 用户相关路由
		userapi.RegisterRoutes(adminGroup, services.UserService, services.UserRepository, services.TokenService)

		// 系统设置相关路由
		setting.RegisterRoutes(adminGroup, services.SettingService)

		// 角色相关路由
		role.RegisterRoutes(adminGroup, services.RoleService)

		// 文件上传相关路由
		upload.RegisterRoutes(adminGroup)

		// 操作日志相关路由
		operate_log.RegisterRoutes(adminGroup, services.OperateLogService)

		// 位置相关路由
		position.RegisterRoutes(adminGroup, services.PositionService)
	}

	// 静态文件服务 - 提供上传文件的访问
	uploadPath := "./uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, 0755)
	}
	r.Static("/uploads", uploadPath)
}