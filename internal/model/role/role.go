package role

import (
	"goadmin/internal/model/permission"
	"goadmin/internal/model/schema"
)

// Role 角色表
type Role struct {
	schema.BaseModel
	Code        string     `gorm:"column:code;size:32;not null;unique;default:''" json:"code"`
	Name        string     `gorm:"column:name;size:50;not null;unique;default:''" json:"name"`
	Description string     `gorm:"column:description;size:200;default:''" json:"description"`
	Status      RoleStatus `gorm:"column:status;default:1;comment:1:active,2:inactive" json:"status"`
	SystemFlag  SystemFlag `gorm:"column:system_flag;default:2;comment:2:非系统,1:系统" json:"system_flag"` //

	Permissions []permission.Permission `gorm:"-" json:"permissions"`
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
	// SystemFlagYes 系统
	SystemFlagYes SystemFlag = iota + 1
	// SystemFlagNo 非系统
	SystemFlagNo
)

type RoleStatus int8

const (
	// RoleStatusActive 激活
	RoleStatusActive RoleStatus = iota + 1
	// RoleStatusInactive 停用
	RoleStatusInactive
)
