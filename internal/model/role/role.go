package role

import (
	"goadmin/internal/model/permission"
	"goadmin/internal/model/schema"
)

// Role 角色表
type Role struct {
	schema.BaseModel
	Code        string     `gorm:"size:32;not null;unique;default:''" json:"code"`
	Name        string     `gorm:"size:50;not null;unique;default:''" json:"name"`
	Description string     `gorm:"size:200;default:''" json:"description"`
	Status      int        `gorm:"default:1;comment:1:active,0:inactive" json:"status"`
	SystemFlag  SystemFlag `gorm:"type:tinyint;default:1;not null" json:"system_flag"` //

	Permissions []permission.Permission `gorm:"many2many:role_permissions;joinForeignKey:RoleCode;joinReferences:PermissionCode"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// SystemFlag 是否系统权限
func (r Role) IsSystem() bool {
	return r.SystemFlag == SystemFlagYes
}

const (
	CodeSuperAdmin = "sup_admin" // 超级管理员角色code
)

type SystemFlag int8

const (
	// SystemFlagNo 非系统
	SystemFlagNo SystemFlag = iota
	// SystemFlagYes 系统
	SystemFlagYes
)
