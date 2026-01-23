package upload

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册上传相关的API路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 创建上传处理器
	handler := NewHandler()

	group := r.Group("/upload")
	{
		// 需要认证的接口
		authGroup := group.Group("/")
		authGroup.Use(middleware.Auth())
		{
			// 上传单个文件
			authGroup.POST("/file", context.Build(handler.UploadFile))
		}
	}
}
