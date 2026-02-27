package setting

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"
	"goadmin/internal/service/setting"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册系统设置相关的API路由
func RegisterRoutes(r *gin.RouterGroup, settingService setting.ServerSettingService) {
	handler := NewHandler(settingService)

	group := r.Group("/setting")
	{
		// 验证码开关配置
		group.GET("/get_settings", context.Build(handler.GetSettings))

		// 需要认证的接口
		authGroup := group.Use(middleware.Auth())
		{
			// 系统设置操作
			authGroup.POST("/set_settings", context.Build(handler.UpdateSettings))

			// 基础配置操作
			authGroup.GET("/get", context.Build(handler.GetByNames))
			authGroup.POST("/set", context.Build(handler.SetByName))

			// 加密配置操作
			authGroup.POST("/encrypted", context.Build(handler.SetEncryptedValue))
			authGroup.GET("/decrypted", context.Build(handler.GetDecryptedValue))
		}

	}
}