package role

import (
	"context"
	"fmt"
	"goadmin/internal/model/permission"
	"goadmin/internal/model/role"
	"goadmin/pkg/db"

	"gorm.io/gorm"
)

// 确保 RolePermissionRepositoryImpl 实现了 RolePermissionRepository 接口
var _ RolePermissionRepository = (*RolePermissionRepositoryImpl)(nil)

// RolePermissionRepositoryImpl 实现 RolePermissionRepository 接口
type RolePermissionRepositoryImpl struct {
	*db.BaseRepository[role.RolePermission]
}

// NewRolePermissionRepository 创建角色权限关联仓储实例
func NewRolePermissionRepository(dbInstance *gorm.DB) RolePermissionRepository {
	return &RolePermissionRepositoryImpl{
		db.NewBaseRepository[role.RolePermission](dbInstance),
	}
}

// NewRolePermissionRepositoryWithDB 使用指定的数据库实例创建角色权限关联仓储实例
func NewRolePermissionRepositoryWithDB() RolePermissionRepository {
	return NewRolePermissionRepository(db.GetDB())
}

// GetPermissionsByRoleCode 根据角色代码获取权限列表
func (r *RolePermissionRepositoryImpl) GetPermissionsByRoleCode(ctx context.Context, roleCode string) ([]string, error) {
	var permissions []string
	err := r.DB().WithContext(ctx).
		Raw("? UNION ?",
			r.DB().Model(&role.RolePermission{}).Select("permission_code").Where("role_code = ?", roleCode),
			r.DB().Model(&permission.Permission{}).Select("code").Where("global_flag = ?", permission.GlobalFlagYes),
		).Pluck("permission_code", &permissions).Error
	return permissions, err
}

// GetPermissionURLsByRoleCode 根据角色代码获取权限列表
func (r *RolePermissionRepositoryImpl) GetPermissionURLsByRoleCode(
	ctx context.Context, roleCode string) ([]string, error) {
	var path []string
	permTable := (permission.Permission{}).TableName()
	roleTable := (role.RolePermission{}).TableName()

	err := r.DB().WithContext(ctx).Model(&role.RolePermission{}).
		Joins(fmt.Sprintf("%s as rp RIGHT JOIN %s as p ON rp.permission_code = p.code", roleTable, permTable)).
		Where("role_code = ? OR p.global_flag = ?", roleCode, permission.GlobalFlagYes).
		Pluck("p.path", &path).Error
	return path, err
}

// GetRolesByPermissionCode 根据权限代码获取角色列表
func (r *RolePermissionRepositoryImpl) GetRolesByPermissionCode(ctx context.Context, permissionCode string) ([]string, error) {
	var roles []string
	err := r.DB().WithContext(ctx).Model(&role.RolePermission{}).
		Where("permission_code = ?", permissionCode).
		Pluck("role_code", &roles).Error
	return roles, err
}

// BatchCreate 批量创建角色权限关联
func (r *RolePermissionRepositoryImpl) BatchCreate(ctx context.Context, rolePermissions []*role.RolePermission) error {
	if len(rolePermissions) == 0 {
		return nil
	}
	return r.DB().WithContext(ctx).Create(&rolePermissions).Error
}

// DeleteByRoleCode 删除角色的所有权限
func (r *RolePermissionRepositoryImpl) DeleteByRoleCode(ctx context.Context, roleCode string) error {
	return r.DB().WithContext(ctx).Where("role_code = ?", roleCode).Delete(&role.RolePermission{}).Error
}

// HasPermission 检查角色是否有特定权限
func (r *RolePermissionRepositoryImpl) HasPermission(ctx context.Context, roleCode string, permissionCode string) (bool, error) {
	var count int64
	err := r.DB().WithContext(ctx).Model(&role.RolePermission{}).
		Where("role_code = ? AND permission_code = ?", roleCode, permissionCode).
		Count(&count).Error
	return count > 0, err
}

// GetPermissionsByRoleCodes 批量获取多个角色的权限列表
func (r *RolePermissionRepositoryImpl) GetPermissionsByRoleCodes(
	ctx context.Context, roleCodes []string) (map[string][]permission.Permission, error) {
	if len(roleCodes) == 0 {
		return make(map[string][]permission.Permission), nil
	}

	cols := []string{"rp.role_code", "p.code", "p.name", "p.path", "p.global_flag", "p.description"}
	var rolePermissions []*role.RoleFullPermission
	err := r.DB().WithContext(ctx).
		Table(role.TableNameRolePermission+" rp").
		Joins("LEFT JOIN "+permission.TableNamePermission+" p ON rp.permission_code = p.code").
		Where("role_code IN ?", roleCodes).
		Select(cols).
		Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}

	// 按角色代码分组权限
	result := make(map[string][]permission.Permission)
	for _, rp := range rolePermissions {
		result[rp.RoleCode] = append(result[rp.RoleCode], rp.Permission)
	}

	return result, nil
}
