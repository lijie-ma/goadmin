package context

import "goadmin/internal/model/role"

type Session interface {
	// GetID 返回会话唯一ID（可对应数据库或token中的用户ID）
	GetID() uint64

	// GetUsername 返回用户名
	GetUsername() string

	// GetStatus 返回用户状态（例如 0=禁用，1=启用）
	GetStatus() int

	// GetRole 返回角色标识（如 "admin", "user"）
	GetRole() *role.Role

	// IsActive 用户状态是否正常
	IsActive() bool
}
