package permission

import "goadmin/internal/model/schema"

const (
	TableNamePermission = "permissions"
)

// Permission 权限表
type Permission struct {
	schema.BaseModel
	Code        string     `gorm:"size:32;not null;unique;default:''" json:"code"` // 权限唯一标识，如 user:create
	Name        string     `gorm:"size:50;not null;unique;default:''" json:"name"`
	Description string     `gorm:"size:200;default:''" json:"description"`
	Path        string     `gorm:"size:200;not null;default:''" json:"path"`           // API路径
	GlobalFlag  GlobalFlag `gorm:"type:tinyint;default:1;not null" json:"global_flag"` //
	Module      string     `gorm:"size:50;not null;default:''" json:"module"`          // 所属模块
}

// TableName 指定表名
func (Permission) TableName() string {
	return TableNamePermission
}

// 是否全局权限
func (p Permission) IsGlobal() bool {
	return p.GlobalFlag == GlobalFlagYes
}

type GlobalFlag int8

const (
	// GlobalFlagNo 非全局
	GlobalFlagNo GlobalFlag = iota
	// GlobalFlagYes 全局
	GlobalFlagYes
)
