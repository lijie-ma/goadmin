package role

import (
	"context"
	"goadmin/internal/model/role"
	"goadmin/pkg/db"

	"gorm.io/gorm"
)

// 确保 RoleRepositoryImpl 实现了 RoleRepository 接口
var _ RoleRepository = (*RoleRepositoryImpl)(nil)

// RoleRepositoryImpl 实现 RoleRepository 接口
type RoleRepositoryImpl struct {
	*db.BaseRepository[role.Role]
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository(dbInstance *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{
		db.NewBaseRepository[role.Role](dbInstance),
	}
}

// NewRoleRepositoryWithDB 使用全局数据库实例创建角色仓储实例
func NewRoleRepositoryWithDB() RoleRepository {
	return NewRoleRepository(db.GetDB())
}

// GetByCode 根据角色代码获取角色
func (r *RoleRepositoryImpl) GetByCode(ctx context.Context, code string) (*role.Role, error) {
	var result role.Role
	err := r.DB().WithContext(ctx).Where("code = ?", code).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// GetByName 根据角色名称获取角色
func (r *RoleRepositoryImpl) GetByName(ctx context.Context, name string) (*role.Role, error) {
	var result role.Role
	err := r.DB().WithContext(ctx).Where("name = ?", name).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ListActive 获取所有激活的角色
func (r *RoleRepositoryImpl) ListActive(ctx context.Context) ([]*role.Role, error) {
	var roles []*role.Role
	err := r.DB().WithContext(ctx).Where("status = ?", 1).Find(&roles).Error
	return roles, err
}

// GetByCodes 根据角色代码列表获取角色
func (r *RoleRepositoryImpl) GetByCodes(ctx context.Context, codes []string) ([]*role.Role, error) {
	if len(codes) == 0 {
		return []*role.Role{}, nil
	}
	var roles []*role.Role
	err := r.DB().WithContext(ctx).Where("code IN ?", codes).Find(&roles).Error
	return roles, err
}

// ListWithPermissions 获取角色列表及其权限
func (r *RoleRepositoryImpl) ListWithPermissions(ctx context.Context, page, pageSize int) ([]*role.Role, int64, error) {
	var roles []*role.Role
	var count int64

	// 计算总数
	err := r.DB().WithContext(ctx).Model(&role.Role{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.DB().WithContext(ctx).
		Preload("Permissions").
		Offset(offset).
		Limit(pageSize).
		Find(&roles).Error

	return roles, count, err
}

// GetWithPermissions 根据ID获取角色及其权限
func (r *RoleRepositoryImpl) GetWithPermissions(ctx context.Context, id uint64) (*role.Role, error) {
	var result role.Role
	err := r.DB().WithContext(ctx).
		Preload("Permissions").
		Where("id = ?", id).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// GetByCodeWithPermissions 根据Code获取角色及其权限
func (r *RoleRepositoryImpl) GetByCodeWithPermissions(ctx context.Context, code string) (*role.Role, error) {
	var result role.Role
	err := r.DB().WithContext(ctx).
		Preload("Permissions").
		Where("code = ?", code).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
