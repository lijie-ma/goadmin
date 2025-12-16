package user

import "goadmin/internal/model/role"

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token"` // 验证码token
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token           string    `json:"token"`         // JWT token
	RefreshToken    string    `json:"refresh_token"` // 刷新token
	ExpiresAt       int64     `json:"expires_at"`
	Username        string    `json:"username"`         // 用户名
	RoleCode        string    `json:"role_code"`        // 角色代码
	Email           string    `json:"email,omitempty"`  // 邮箱
	Role            role.Role `json:"role"`             // 角色
	PermissionCodes []string  `json:"permission_codes"` // 权限
}

// ChangePasswordRequest 修改密码请求参数
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
