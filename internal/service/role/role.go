package role

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/role"
	"goadmin/internal/model/schema"
	rolerepo "goadmin/internal/repository/role"
	"goadmin/pkg/db"
	"goadmin/pkg/util"

	"goadmin/config"
)

// RoleService 角色服务接口
type RoleService interface {
	// GetRoleByID 根据ID获取角色
	GetRoleByID(ctx *context.Context, id uint64) (*role.Role, error)

	// GetRoleByCode 根据Code获取角色
	GetRoleByCode(ctx *context.Context, code string) (*role.Role, error)

	// ListRoles 获取角色列表
	ListRoles(ctx *context.Context, req *schema.PageRequest) ([]*role.Role, int64, error)

	// ListActiveRoles 获取所有激活的角色
	ListActiveRoles(ctx *context.Context) ([]*role.Role, error)

	// CreateRole 创建角色
	CreateRole(ctx *context.Context, roleModel *role.CreateRequest) error

	// UpdateRole 更新角色
	UpdateRole(ctx *context.Context, roleModel *role.UpdateRequest) error

	// DeleteRole 删除角色
	DeleteRole(ctx *context.Context, id uint64) error

	// GetRoleWithPermissions 获取角色及其权限
	GetRoleWithPermissions(ctx *context.Context, id uint64) (*role.Role, error)

	// GetRolePermissions 获取角色的权限列表
	GetRolePermissions(ctx *context.Context, roleCode string) ([]string, error)

	// AssignPermissions 分配权限给角色
	AssignPermissions(ctx *context.Context, roleCode string, permissionCodes []string) error

	// HasPermission 检查角色是否有特定权限
	HasPermission(ctx *context.Context, roleCode string, permissionCode string) (bool, error)

	// ListAllPermissions 获取所有权限列表
	ListAllPermissions(ctx *context.Context) ([]map[string]interface{}, error)
}

// roleService 角色服务实现
type roleService struct {
	roleRepo           rolerepo.RoleRepository
	rolePermissionRepo rolerepo.RolePermissionRepository
	cfg                *config.Config
}

// NewRoleService 创建角色服务实例
func NewRoleService() RoleService {
	return &roleService{
		cfg:                config.Get(),
		roleRepo:           rolerepo.NewRoleRepositoryWithDB(),
		rolePermissionRepo: rolerepo.NewRolePermissionRepositoryWithDB(),
	}
}

func (*roleService) logPrefix() string {
	return "role-service"
}

// GetRoleByID 根据ID获取角色
func (s *roleService) GetRoleByID(ctx *context.Context, id uint64) (*role.Role, error) {
	rs, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色失败: %d %v", s.logPrefix(), id, err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if rs == nil {
		return nil, i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
	}
	return rs, nil
}

// GetRoleByCode 根据Code获取角色
func (s *roleService) GetRoleByCode(ctx *context.Context, code string) (*role.Role, error) {
	rs, err := s.roleRepo.GetByCode(ctx, code)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色失败: %s %v", s.logPrefix(), code, err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if rs == nil {
		return nil, i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
	}
	return rs, nil
}

// ListRoles 获取角色列表
func (s *roleService) ListRoles(ctx *context.Context, req *schema.PageRequest) ([]*role.Role, int64, error) {
	list, total, err := s.roleRepo.List(ctx, req.Page, req.PageSize, db.Order[role.Role](req.OrderBy))
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色列表失败: %v", s.logPrefix(), err)
		return nil, 0, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if total == 0 {
		return []*role.Role{}, 0, nil
	}
	var (
		roleCodes []string
	)
	for _, r := range list {
		roleCodes = append(roleCodes, r.Code)
	}
	permissions, err := s.rolePermissionRepo.GetPermissionsByRoleCodes(ctx, roleCodes)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色权限失败: %v", s.logPrefix(), err)
		return nil, 0, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	for _, r := range list {
		r.Permissions = permissions[r.Code]
	}
	return list, total, nil
}

// ListActiveRoles 获取所有激活的角色
func (s *roleService) ListActiveRoles(ctx *context.Context) ([]*role.Role, error) {
	return s.roleRepo.ListActive(ctx)
}

// CreateRole 创建角色
func (s *roleService) CreateRole(ctx *context.Context, roleModel *role.CreateRequest) error {
	var (
		code         string
		err          error
		existingRole *role.Role
	)
	for i := range 3 {
		code, err = util.GenerateRandomString(8)
		if err != nil {
			ctx.Logger.Errorf("%s 生成角色代码失败: times=%d %v", s.logPrefix(), i, err)
			continue
		}
		// 检查角色代码是否已存在
		existingRole, err = s.roleRepo.GetByCode(ctx, roleModel.Code)
		if err != nil {
			ctx.Logger.Errorf("%s 检查角色代码是否存在失败: %s %v", s.logPrefix(), roleModel.Code, err)
			continue
		}
		if existingRole == nil {
			break
		}
	}
	if err != nil {
		ctx.Logger.Errorf("%s 检查角色代码是否存在失败: %s %v", s.logPrefix(), roleModel.Code, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if existingRole != nil {
		return i18n.E(
			ctx.Context, "common.HadExist", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
	}

	// 检查角色名称是否已存在
	existingRole, err = s.roleRepo.GetByName(ctx, roleModel.Name)
	if err != nil {
		ctx.Logger.Errorf("%s 检查角色名称是否存在失败: %s %v", s.logPrefix(), roleModel.Name, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if existingRole != nil {
		return i18n.E(
			ctx.Context, "common.HadExist", map[string]any{"item": roleModel.Name})
	}

	// 创建角色
	err = s.roleRepo.Create(ctx, &role.Role{
		Code:        code,
		Name:        roleModel.Name,
		Description: roleModel.Description,
		Status:      roleModel.Status,
		SystemFlag:  role.SystemFlagNo,
	})
	if err != nil {
		ctx.Logger.Errorf("%s 创建角色失败: %v", s.logPrefix(), err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	return nil
}

// UpdateRole 更新角色
func (s *roleService) UpdateRole(ctx *context.Context, roleModel *role.UpdateRequest) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByID(ctx, roleModel.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色失败: %d %v", s.logPrefix(), roleModel.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if existingRole == nil {
		return i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
	}

	if existingRole.IsSystem() {
		return i18n.E(ctx.Context, "common.PermissionDeny", nil)
	}

	// 如果更改了名称，检查新名称是否已存在
	if existingRole.Name != roleModel.Name {
		nameExists, err := s.roleRepo.GetByName(ctx, roleModel.Name)
		if err != nil {
			ctx.Logger.Errorf("%s 检查角色名称是否存在失败: %s %v", s.logPrefix(), roleModel.Name, err)
			return i18n.E(ctx.Context, "common.RepositoryErr", nil)
		}

		if nameExists != nil {
			return i18n.E(ctx.Context, "common.HadExist", map[string]any{"item": roleModel.Name})
		}
	}
	existingRole.Name = roleModel.Name
	// existingRole.Code = roleModel.Code  // code 禁止修改
	existingRole.Description = roleModel.Description
	existingRole.Status = roleModel.Status
	existingRole.MTime = util.Now()

	// 更新角色
	err = s.roleRepo.Update(ctx, existingRole)
	if err != nil {
		ctx.Logger.Errorf("%s 更新角色失败: %v", s.logPrefix(), err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	return nil
}

// DeleteRole 删除角色
func (s *roleService) DeleteRole(ctx *context.Context, id uint64) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色失败: %d %v", s.logPrefix(), id, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if existingRole == nil {
		return i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
	}

	if existingRole.IsSystem() {
		return i18n.E(ctx.Context, "common.PermissionDeny", nil)
	}

	// 删除角色权限关联
	err = s.rolePermissionRepo.DeleteByRoleCode(ctx, existingRole.Code)
	if err != nil {
		ctx.Logger.Errorf("%s 删除角色权限关联失败: %s %v", s.logPrefix(), existingRole.Code, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	// 删除角色
	err = s.roleRepo.Delete(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("%s 删除角色权限关联失败: %s %v", s.logPrefix(), existingRole.Code, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	return nil
}

// GetRoleWithPermissions 获取角色及其权限
func (s *roleService) GetRoleWithPermissions(ctx *context.Context, id uint64) (*role.Role, error) {
	return s.roleRepo.GetWithPermissions(ctx, id)
}

// GetRolePermissions 获取角色的权限列表
func (s *roleService) GetRolePermissions(ctx *context.Context, roleCode string) ([]string, error) {
	return s.rolePermissionRepo.GetPermissionsByRoleCode(ctx, roleCode)
}

// AssignPermissions 分配权限给角色
func (s *roleService) AssignPermissions(ctx *context.Context, roleCode string, permissionCodes []string) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByCode(ctx, roleCode)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色失败: %s %v", s.logPrefix(), roleCode, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if existingRole == nil {
		return i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
	}

	// 删除当前角色的所有权限
	err = s.rolePermissionRepo.DeleteByRoleCode(ctx, roleCode)
	if err != nil {
		ctx.Logger.Errorf("%s 删除角色权限关联失败: %s %v", s.logPrefix(), roleCode, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	// 如果没有新权限，直接返回
	if len(permissionCodes) == 0 {
		return nil
	}

	// 创建新的角色权限关联
	rolePermissions := make([]*role.RolePermission, 0, len(permissionCodes))
	for _, permCode := range permissionCodes {
		rolePermissions = append(rolePermissions, &role.RolePermission{
			RoleCode:       roleCode,
			PermissionCode: permCode,
		})
	}

	// 批量创建角色权限关联
	err = s.rolePermissionRepo.BatchCreate(ctx, rolePermissions)
	if err != nil {
		ctx.Logger.Errorf("%s 批量创建角色权限关联失败: %s %v", s.logPrefix(), roleCode, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	return nil
}

// HasPermission 检查角色是否有特定权限
func (s *roleService) HasPermission(ctx *context.Context, roleCode string, permissionCode string) (bool, error) {
	return s.rolePermissionRepo.HasPermission(ctx, roleCode, permissionCode)
}

// ListAllPermissions 获取所有权限列表
func (s *roleService) ListAllPermissions(ctx *context.Context) ([]map[string]interface{}, error) {
	permissions, err := s.rolePermissionRepo.GetAllPermissions(ctx)
	if err != nil {
		ctx.Logger.Errorf("%s 获取所有权限列表失败: %v", s.logPrefix(), err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	// 将权限按模块分组
	moduleMap := make(map[string][]map[string]interface{})
	for _, perm := range permissions {
		permMap := map[string]interface{}{
			"code":        perm.Code,
			"name":        perm.Name,
			"description": perm.Description,
			"module":      perm.Module,
		}
		moduleMap[perm.Module] = append(moduleMap[perm.Module], permMap)
	}

	// 构建返回结果
	var result []map[string]interface{}
	for module, perms := range moduleMap {
		result = append(result, map[string]interface{}{
			"module":      module,
			"permissions": perms,
		})
	}

	return result, nil
}
