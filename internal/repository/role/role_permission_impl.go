package role

import (
	"context"
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
func (r *RolePermissionRepositoryImpl) GetPermissionsByRoleCode(
	ctx context.Context, roleCode string, containPublic ...bool) ([]string, error) {
	var (
		permissions []string
		err         error
		db          *gorm.DB
	)
	db = r.DB().WithContext(ctx)
	if len(containPublic) == 0 || containPublic[0] {
		db = db.Raw("? UNION ?",
			r.DB().Model(&role.RolePermission{}).Select("permission_code").Where("role_code = ?", roleCode),
			r.DB().Model(&permission.Permission{}).Select("code").Where("global_flag = ?", permission.GlobalFlagYes),
		)
	} else {
		db = db.Model(&role.RolePermission{}).Select("permission_code").Where("role_code = ?", roleCode)
	}
	err = db.Pluck("permission_code", &permissions).Error
	return permissions, err
}

// GetPermissionURLsByRoleCode 根据角色代码获取权限列表
func (r *RolePermissionRepositoryImpl) GetPermissionURLsByRoleCode(
	ctx context.Context, roleCode string) ([]string, error) {
	var path []string
	permTable := (permission.Permission{}).TableName()
	rolePermTable := (role.RolePermission{}).TableName()

	err := r.DB().WithContext(ctx).Table(permTable+" as p").
		Joins("LEFT JOIN "+rolePermTable+" AS rp ON rp.permission_code = p.code").
		Where("rp.role_code = ? OR p.global_flag = ?", roleCode, permission.GlobalFlagYes).
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

// GetAllPermissions 获取所有权限列表
func (r *RolePermissionRepositoryImpl) GetAllPermissions(
	ctx context.Context, containPublic ...bool) ([]permission.Permission, error) {
	var (
		permissions []permission.Permission
		err         error
		db          *gorm.DB
	)
	db = r.DB().WithContext(ctx).Model(&permission.Permission{})
	if len(containPublic) > 0 && !containPublic[0] {
		db = db.Where("global_flag != ?", permission.GlobalFlagYes)
	}
	err = db.Find(&permissions).Error
	return permissions, err
}

// HasAccessURL 检查角色是否有访问指定URL的权限
func (r *RolePermissionRepositoryImpl) HasAccessURL(
	ctx context.Context, roleCode, accessURL string) (bool, error) {
	db := r.DB().WithContext(ctx)
	permTable := (permission.Permission{}).TableName()
	rolePermTable := (role.RolePermission{}).TableName()

	var cnt int64
	// 检查角色权限
	err := db.Table(rolePermTable+" AS rp").
		Joins("JOIN "+permTable+" AS p ON rp.permission_code = p.code").
		Where("rp.role_code = ? AND p.path = ?", roleCode, accessURL).
		Count(&cnt).Error
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}

	// 检查全局权限
	err = db.Table(permTable).
		Where("path = ? AND global_flag = ?", accessURL, permission.GlobalFlagYes).
		Count(&cnt).Error
	return cnt > 0, err
}
