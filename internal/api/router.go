package api

import (
	"goadmin/internal/api/admin/v1/captcha"
	"goadmin/internal/api/admin/v1/role"
	"goadmin/internal/api/admin/v1/setting"
	"goadmin/internal/api/admin/v1/user"
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
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
	adminHandler(r)
}

func adminHandler(r *gin.Engine) {
	adminGroup := r.Group("/admin/v1")
	{
		// 验证码路由
		captchaRouter := adminGroup.Group("/captcha")
		{
			captchaRouter.GET("/generate", context.Build(captcha.GenerateHandler))
			captchaRouter.POST("/check", context.Build(captcha.CheckHandler))
		}

		// 用户相关路由
		user.RegisterRoutes(adminGroup)

		// 系统设置相关路由
		setting.RegisterRoutes(adminGroup)

		// 角色相关路由
		role.RegisterRoutes(adminGroup)
	}
}
