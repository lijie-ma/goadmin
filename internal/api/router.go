package api

import (
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
	api := r.Group("/api")
	{
		// 用户相关API
		api.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "获取用户列表",
			})
		})

		// 角色相关API
		api.GET("/roles", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "获取角色列表",
			})
		})

		// 权限相关API
		api.GET("/permissions", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "获取权限列表",
			})
		})
	}
}
