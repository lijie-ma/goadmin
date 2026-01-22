package user

import (
	"goadmin/internal/model/schema"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token"` // 验证码token
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string `json:"token"`         // JWT token
	RefreshToken string `json:"refresh_token"` // 刷新token
	ExpiresAt    int64  `json:"expires_at"`
}

// ChangePasswordRequest 修改密码请求参数
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// ListRequest 用户列表请求参数
type ListRequest struct {
	schema.PageRequest
	Keyword string `form:"keyword"` // 搜索关键词（用户名或邮箱）
}

// CreateUserRequest 创建用户请求参数
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Password string `json:"password" binding:"required,min=6,max=32"` // 密码
	Email    string `json:"email" binding:"omitempty,email"`          // 邮箱
	RoleCode string `json:"role_code" binding:"required"`             // 角色代码
	Status   int    `json:"status" binding:"omitempty,min=0,max=1"`   // 状态：0-禁用，1-启用
}

// UpdateUserRequest 更新用户请求参数
type UpdateUserRequest struct {
	ID       uint64 `json:"id" binding:"required"`                     // 用户ID
	Username string `json:"username" binding:"omitempty,min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"omitempty,email"`           // 邮箱
	RoleCode string `json:"role_code" binding:"omitempty"`             // 角色代码
	Status   int    `json:"status" binding:"omitempty,min=0,max=1"`    // 状态：0-禁用，1-启用
}
