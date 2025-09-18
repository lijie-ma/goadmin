package user

import "goadmin/internal/model/schema"

// User 用户表
type User struct {
	schema.BaseModel
	Username string     `gorm:"size:50;not null;unique;default:''" json:"username"`
	Password string     `gorm:"size:100;not null;default:''" json:"password"`
	Email    string     `gorm:"size:100;unique;default:''" json:"email"`
	Status   UserStatus `gorm:"type:int;default:1;comment:0:inactive,1:active,2:locked,3:deleted" json:"status"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
