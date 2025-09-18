package permission

import "goadmin/internal/model/schema"

// Permission 权限表
type Permission struct {
	schema.BaseModel
	Code        string `gorm:"size:32;not null;unique;default:''" json:"code"` // 权限唯一标识，如 user:create
	Name        string `gorm:"size:50;not null;unique;default:''" json:"name"`
	Description string `gorm:"size:200;default:''" json:"description"`
	Path        string `gorm:"size:200;not null;default:''" json:"path"`     // API路径
	Method      string `gorm:"size:20;not null;default:'GET'" json:"method"` // HTTP方法
	Module      string `gorm:"size:50;not null;default:''" json:"module"`    // 所属模块
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
