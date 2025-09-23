package setting

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"
	"goadmin/internal/service/setting"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册系统设置相关的API路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 创建依赖服务
	settingService := setting.NewServerSettingService()
	handler := NewHandler(settingService)

	group := r.Group("/setting")
	{
		// 需要认证的接口
		group.Use(middleware.Auth())
		{
			// 基础配置操作
			group.GET("/:name", context.Build(handler.GetByName))
			group.PUT("/:name", context.Build(handler.SetByName))
			group.POST("/batch", context.Build(handler.BatchGetValues))

			// 验证码开关配置
			group.GET("/captcha-switch", context.Build(handler.GetCaptchaSwitch))
			group.POST("/captcha-switch", context.Build(handler.SetCaptchaSwitch))
		}
	}
}
