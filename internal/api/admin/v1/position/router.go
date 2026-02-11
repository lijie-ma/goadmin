package position

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"
	positionSrv "goadmin/internal/service/position"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册位置相关的API路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 创建依赖服务
	positionService := positionSrv.NewPositionService_legacy()
	handler := NewHandler(positionService)

	group := r.Group("/position")
	{
		// 需要认证的接口
		authGroup := group.Group("/")
		authGroup.Use(middleware.Auth())
		{
			authGroup.GET("/list", context.Build(handler.ListPositions))
			authGroup.GET("/get", context.Build(handler.GetPosition))
			authGroup.POST("/create", context.Build(handler.CreatePosition))
			authGroup.POST("/update", context.Build(handler.UpdatePosition))
			authGroup.POST("/delete", context.Build(handler.DeletePosition))
		}
	}
}
