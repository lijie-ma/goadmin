package role

import "goadmin/internal/model/schema"

// Role 角色表
type Role struct {
	schema.BaseModel
	Code        string `gorm:"size:32;not null;unique;default:''" json:"code"`
	Name        string `gorm:"size:50;not null;unique;default:''" json:"name"`
	Description string `gorm:"size:200;default:''" json:"description"`
	Status      int    `gorm:"default:1;comment:1:active,0:inactive" json:"status"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
