package role

import (
	"context"
	"goadmin/internal/model/role"
	"goadmin/pkg/db"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	db.Repository[role.Role]

	// GetByCode 根据角色代码获取角色
	GetByCode(ctx context.Context, code string) (*role.Role, error)

	// GetByName 根据角色名称获取角色
	GetByName(ctx context.Context, name string) (*role.Role, error)

	// ListActive 获取所有激活的角色
	ListActive(ctx context.Context) ([]*role.Role, error)

	// GetByCodes 根据角色代码列表获取角色
	GetByCodes(ctx context.Context, codes []string) ([]*role.Role, error)

	// ListWithPermissions 获取角色列表及其权限
	ListWithPermissions(ctx context.Context, page, pageSize int) ([]*role.Role, int64, error)

	// GetWithPermissions 根据ID获取角色及其权限
	GetWithPermissions(ctx context.Context, id uint64) (*role.Role, error)

	// GetByCodeWithPermissions 根据Code获取角色及其权限
	GetByCodeWithPermissions(ctx context.Context, code string) (*role.Role, error)
}
