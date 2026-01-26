package operate_log

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"
	operatelog "goadmin/internal/service/operate_log"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册操作日志相关的API路由
func RegisterRoutes(r *gin.RouterGroup) {

	// 创建依赖服务
	logService := operatelog.NewOperateLogService()
	handler := NewHandler(logService)

	group := r.Group("/operate_log")
	{
		// 需要认证的接口
		authGroup := group.Group("/")

		authGroup.Use(middleware.Auth())
		{
			authGroup.GET("/list", context.Build(handler.ListOperateLogs))
		}
	}
}
