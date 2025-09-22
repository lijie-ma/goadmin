package role

import (
	"context"
	"goadmin/internal/model/role"
	"goadmin/pkg/db"
)

// RolePermissionRepository 角色权限关联仓储接口
type RolePermissionRepository interface {
	db.Repository[role.RolePermission]

	// GetPermissionsByRoleCode 根据角色代码获取权限列表
	GetPermissionsByRoleCode(ctx context.Context, roleCode string) ([]string, error)

	// GetPermissionURLsByRoleCode 根据角色代码获取权限列表
	GetPermissionURLsByRoleCode(ctx context.Context, roleCode string) ([]string, error)

	// GetRolesByPermissionCode 根据权限代码获取角色列表
	GetRolesByPermissionCode(ctx context.Context, permissionCode string) ([]string, error)

	// BatchCreate 批量创建角色权限关联
	BatchCreate(ctx context.Context, rolePermissions []*role.RolePermission) error

	// DeleteByRoleCode 删除角色的所有权限
	DeleteByRoleCode(ctx context.Context, roleCode string) error

	// HasPermission 检查角色是否有特定权限
	HasPermission(ctx context.Context, roleCode string, permissionCode string) (bool, error)

	// GetPermissionsByRoleCodes 批量获取多个角色的权限列表
	GetPermissionsByRoleCodes(ctx context.Context, roleCodes []string) (map[string][]string, error)
}
