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
			group.GET("/list", context.Build(handler.ListRoles))
			group.GET("/active", context.Build(handler.ListActiveRoles))
			group.POST("/get", context.Build(handler.GetRole))
			group.POST("/create", context.Build(handler.CreateRole))
			group.POST("/update", context.Build(handler.UpdateRole))
			group.POST("/delete", context.Build(handler.DeleteRole))

			// 角色权限管理
			group.POST("/permissions/get", context.Build(handler.GetRolePermissions))
			group.POST("/permissions/assign", context.Build(handler.AssignPermissions))
		}
	}
}
