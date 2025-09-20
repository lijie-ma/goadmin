package api

import (
	"goadmin/internal/api/admin/v1/captcha"
	"goadmin/internal/context"
	"goadmin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	r.Use(
		middleware.Trace(),
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
		captchaRouter := adminGroup.Group("/captcha")
		{
			captchaRouter.GET("/generate", context.Build(captcha.GenerateHandler))
			captchaRouter.POST("/check", context.Build(captcha.CheckHandler))
		}

	}
}
