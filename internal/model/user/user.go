package user

import (
	"goadmin/internal/model/role"
	"goadmin/internal/model/schema"
)

// User 用户表
type User struct {
	schema.BaseModel
	Username string     `gorm:"size:50;not null;unique;default:''" json:"username"`
	Password string     `gorm:"size:100;not null;default:''" json:"-"`
	Email    string     `gorm:"size:100;unique;default:''" json:"email"`
	Status   UserStatus `gorm:"type:int;default:1;comment:0:inactive,1:active,2:locked,3:deleted" json:"status"`
	RoleCode string     `gorm:"size:32;not null;index:idx_user_role;default:''" json:"role_code"`
	// 用户角色关联表
	Role role.Role `gorm:"foreignKey:RoleCode;references:Code" json:"role"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) IsSuperAdmin() bool {
	return u.RoleCode == role.CodeSuperAdmin
}

// GetID 实现 Session 接口
func (u *User) GetID() uint64 {
	return u.ID
}

// GetUsername 实现 Session 接口
func (u *User) GetUsername() string {
	return u.Username
}

// GetStatus 实现 Session 接口
func (u *User) GetStatus() int {
	return int(u.Status)
}

// GetRole 实现 Session 接口
func (u *User) GetRole() *role.Role {
	return &u.Role
}
