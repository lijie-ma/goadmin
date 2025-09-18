package user

import "goadmin/internal/model/schema"

// UserRole 用户角色关联表
type UserRole struct {
	schema.BaseModel
	UserID   uint64 `gorm:"not null;index:idx_user_role;default:0" json:"user_id"`
	RoleCode string `gorm:"size:32;not null;index:idx_user_role;default:''" json:"role_code"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
