package role

import (
	"fmt"
	"goadmin/internal/context"
	"goadmin/internal/model/role"
	"goadmin/internal/model/schema"
	rolerepo "goadmin/internal/repository/role"
	"goadmin/pkg/db"

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
	CreateRole(ctx *context.Context, roleModel *role.Role) error

	// UpdateRole 更新角色
	UpdateRole(ctx *context.Context, roleModel *role.Role) error

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
	return s.roleRepo.GetByID(ctx, id)
}

// GetRoleByCode 根据Code获取角色
func (s *roleService) GetRoleByCode(ctx *context.Context, code string) (*role.Role, error) {
	return s.roleRepo.GetByCode(ctx, code)
}

// ListRoles 获取角色列表
func (s *roleService) ListRoles(ctx *context.Context, req *schema.PageRequest) ([]*role.Role, int64, error) {
	list, total, err := s.roleRepo.List(ctx, req.Page, req.PageSize, db.Order[role.Role](req.OrderBy))
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色列表失败: %v", s.logPrefix(), err)
		return nil, 0, err
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
		return nil, 0, err
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
func (s *roleService) CreateRole(ctx *context.Context, roleModel *role.Role) error {
	// 检查角色代码是否已存在
	existingRole, err := s.roleRepo.GetByCode(ctx, roleModel.Code)
	if err != nil {
		ctx.Logger.Errorf("%s 检查角色代码是否存在失败: %s %v", s.logPrefix(), roleModel.Code, err)
		return err
	}

	if existingRole != nil {
		return fmt.Errorf("角色代码 %s 已存在", roleModel.Code)
	}

	// 检查角色名称是否已存在
	existingRole, err = s.roleRepo.GetByName(ctx, roleModel.Name)
	if err != nil {
		ctx.Logger.Errorf("%s 检查角色名称是否存在失败: %s %v", s.logPrefix(), roleModel.Name, err)
		return err
	}

	if existingRole != nil {
		return fmt.Errorf("角色名称 %s 已存在", roleModel.Name)
	}

	// 创建角色
	return s.roleRepo.Create(ctx, roleModel)
}

// UpdateRole 更新角色
func (s *roleService) UpdateRole(ctx *context.Context, roleModel *role.Role) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByID(ctx, roleModel.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色失败: %d %v", s.logPrefix(), roleModel.ID, err)
		return err
	}

	if existingRole == nil {
		return fmt.Errorf("角色不存在")
	}

	if existingRole.IsSystem() {
		return fmt.Errorf("系统角色不能被修改")
	}

	// 如果更改了代码，检查新代码是否已存在
	if existingRole.Code != roleModel.Code {
		codeExists, err := s.roleRepo.GetByCode(ctx, roleModel.Code)
		if err != nil {
			ctx.Logger.Errorf("%s 检查角色代码是否存在失败: %s %v", s.logPrefix(), roleModel.Code, err)
			return err
		}

		if codeExists != nil {
			return fmt.Errorf("角色代码 %s 已存在", roleModel.Code)
		}
	}

	// 如果更改了名称，检查新名称是否已存在
	if existingRole.Name != roleModel.Name {
		nameExists, err := s.roleRepo.GetByName(ctx, roleModel.Name)
		if err != nil {
			ctx.Logger.Errorf("%s 检查角色名称是否存在失败: %s %v", s.logPrefix(), roleModel.Name, err)
			return err
		}

		if nameExists != nil {
			return fmt.Errorf("角色名称 %s 已存在", roleModel.Name)
		}
	}

	// 更新角色
	return s.roleRepo.Update(ctx, roleModel)
}

// DeleteRole 删除角色
func (s *roleService) DeleteRole(ctx *context.Context, id uint64) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("%s 获取角色失败: %d %v", s.logPrefix(), id, err)
		return err
	}

	if existingRole == nil {
		return fmt.Errorf("角色不存在")
	}

	if existingRole.IsSystem() {
		return fmt.Errorf("不能删除系统内置角色")
	}

	// 删除角色权限关联
	err = s.rolePermissionRepo.DeleteByRoleCode(ctx, existingRole.Code)
	if err != nil {
		ctx.Logger.Errorf("%s 删除角色权限关联失败: %s %v", s.logPrefix(), existingRole.Code, err)
		return err
	}

	// 删除角色
	return s.roleRepo.Delete(ctx, id)
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
		return err
	}

	if existingRole == nil {
		return fmt.Errorf("角色不存在")
	}

	// 删除当前角色的所有权限
	err = s.rolePermissionRepo.DeleteByRoleCode(ctx, roleCode)
	if err != nil {
		ctx.Logger.Errorf("%s 删除角色权限关联失败: %s %v", s.logPrefix(), roleCode, err)
		return err
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
	return s.rolePermissionRepo.BatchCreate(ctx, rolePermissions)
}

// HasPermission 检查角色是否有特定权限
func (s *roleService) HasPermission(ctx *context.Context, roleCode string, permissionCode string) (bool, error) {
	return s.rolePermissionRepo.HasPermission(ctx, roleCode, permissionCode)
}
