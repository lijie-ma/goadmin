package role

import (
	"goadmin/internal/model/permission"
	"goadmin/internal/model/schema"
)

// Role 角色表
type Role struct {
	schema.BaseModel
	Code        string `gorm:"size:32;not null;unique;default:''" json:"code"`
	Name        string `gorm:"size:50;not null;unique;default:''" json:"name"`
	Description string `gorm:"size:200;default:''" json:"description"`
	Status      int    `gorm:"default:1;comment:1:active,0:inactive" json:"status"`

	Permissions []permission.Permission `gorm:"many2many:role_permissions;foreignKey:Code;joinForeignKey:RoleCode;References:PermissionCode;JoinReferences:Code" json:"permissions"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

const (
	CodeSuperAdmin = "sup_admin" // 超级管理员角色code
)
