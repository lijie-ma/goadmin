package role

import "goadmin/internal/model/schema"

// RolePermission 角色权限关联表
type RolePermission struct {
	schema.BaseModel
	RoleCode       string `gorm:"size:32;not null;index:idx_role_permission;default:''" json:"role_code"`
	PermissionCode string `gorm:"size:32;not null;index:idx_role_permission;default:''" json:"permission_code"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
