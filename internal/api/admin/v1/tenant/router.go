package tenant

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"
	tenantSrv "goadmin/internal/service/tenant"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册租户相关的API路由
func RegisterRoutes(r *gin.RouterGroup, tenantService tenantSrv.TenantService) {
	handler := NewHandler(tenantService)

	group := r.Group("/tenant")
	{
		// 需要认证的接口
		authGroup := group.Group("/")
		authGroup.Use(middleware.Auth())
		{
			authGroup.GET("/list", context.Build(handler.ListTenants))
			authGroup.GET("/get", context.Build(handler.GetTenant))
			authGroup.POST("/create", context.Build(handler.CreateTenant))
			authGroup.POST("/update", context.Build(handler.UpdateTenant))
			authGroup.POST("/delete", context.Build(handler.DeleteTenant))
		}
	}
}
