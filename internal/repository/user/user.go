package user

import (
	"context"
	"goadmin/internal/model/user"
	"goadmin/pkg/db"
)

// UserRepository 定义用户仓储接口
type UserRepository interface {
	db.Repository[user.User]

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*user.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*user.User, error)

	// UpdateStatus 更新用户状态
	UpdateStatus(ctx context.Context, id uint64, status user.UserStatus) error

	// UpdatePassword 更新用户密码
	UpdatePassword(ctx context.Context, id uint64, password string) error

	// IsUsernameExists 检查用户名是否存在
	IsUsernameExists(ctx context.Context, username string, excludeID ...uint64) (bool, error)

	// IsEmailExists 检查邮箱是否存在
	IsEmailExists(ctx context.Context, email string, excludeID ...uint64) (bool, error)

	// GetUsersByRoleCode 获取指定角色的所有用户
	GetUsersByRoleCode(ctx context.Context, roleCode string) ([]*user.User, error)

	// List 获取用户列表（重写以排除已删除用户）
	PageList(ctx context.Context, req *user.ListRequest) ([]*user.User, int64, error)
}
