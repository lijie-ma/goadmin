package user

import (
	"context"
	"errors"
	"goadmin/internal/model/user"
	"goadmin/pkg/db"
	"time"
)

// 确保UserRepositoryImpl实现了UserRepository接口
var _ UserRepository = (*UserRepositoryImpl)(nil)

// UserRepositoryImpl 实现UserRepository接口
type UserRepositoryImpl struct {
	*db.BaseRepository[user.User]
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		db.NewBaseRepository[user.User](db.GetDB()),
	}
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uint64) (*user.User, error) {
	var u user.User
	err := r.DB().WithContext(ctx).
		Where("id = ? AND status != ?", id, user.UserStatusDeleted).
		Preload("Role").First(&u).Error
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.DB().WithContext(ctx).
		Where("username = ? AND status != ?", username, user.UserStatusDeleted).
		Preload("Role").First(&u).Error
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.DB().WithContext(ctx).Where("email = ? AND status != ?", email, user.UserStatusDeleted).First(&u).Error
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// UpdateStatus 更新用户状态
func (r *UserRepositoryImpl) UpdateStatus(ctx context.Context, id uint64, status user.UserStatus) error {
	return r.DB().WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
			"mtime":  time.Now(),
		}).Error
}

// UpdatePassword 更新用户密码
func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, id uint64, password string) error {
	return r.DB().WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password": password,
			"mtime":    time.Now(),
		}).Error
}

// Delete 删除用户（逻辑删除）
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.UpdateStatus(ctx, id, user.UserStatusDeleted)
}

// IsUsernameExists 检查用户名是否存在
func (r *UserRepositoryImpl) IsUsernameExists(ctx context.Context, username string, excludeID ...uint64) (bool, error) {
	var count int64
	query := r.DB().WithContext(ctx).Model(&user.User{}).
		Where("username = ? AND status != ?", username, user.UserStatusDeleted)

	// 排除指定ID
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// IsEmailExists 检查邮箱是否存在
func (r *UserRepositoryImpl) IsEmailExists(ctx context.Context, email string, excludeID ...uint64) (bool, error) {
	var count int64
	query := r.DB().WithContext(ctx).Model(&user.User{}).
		Where("email = ? AND status != ?", email, user.UserStatusDeleted)

	// 排除指定ID
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetUsersByRoleCode 获取指定角色的所有用户
func (r *UserRepositoryImpl) GetUsersByRoleCode(ctx context.Context, roleCode string) ([]*user.User, error) {
	var users []*user.User
	err := r.DB().WithContext(ctx).
		Where("role_code = ? AND status != ?", roleCode, user.UserStatusDeleted).
		Find(&users).Error
	return users, err
}
