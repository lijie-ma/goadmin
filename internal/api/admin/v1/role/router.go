package role

import (
	"goadmin/internal/context"
	"goadmin/internal/middleware"
	rolesrv "goadmin/internal/service/role"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册角色相关的API路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 创建依赖服务
	roleService := rolesrv.NewRoleService()
	handler := NewHandler(roleService)

	group := r.Group("/role")
	{
		// 需要认证的接口
		group.Use(middleware.Auth())
		{
			// 角色管理
			group.GET("", context.Build(handler.ListRoles))
			group.GET("/active", context.Build(handler.ListActiveRoles))
			group.GET("/:id", context.Build(handler.GetRole))
			group.POST("", context.Build(handler.CreateRole))
			group.POST("/:id", context.Build(handler.UpdateRole))
			group.DELETE("/:id", context.Build(handler.DeleteRole))

			// 角色权限管理
			group.GET("/permissions/:code", context.Build(handler.GetRolePermissions))
			group.POST("/permissions", context.Build(handler.AssignPermissions))
		}
	}
}
