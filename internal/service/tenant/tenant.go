package tenant

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/schema"
	"goadmin/internal/model/tenant"
	tenantrepo "goadmin/internal/repository/tenant"
	"goadmin/internal/service/operate_log"
)

// TenantService 租户服务接口
type TenantService interface {
	// GetTenantByID 根据ID获取租户
	GetTenantByID(ctx *context.Context, id uint64) (*tenant.Tenant, error)

	// ListTenants 获取租户列表
	ListTenants(ctx *context.Context, req *tenant.ListRequest) ([]*tenant.Tenant, int64, error)

	// CreateTenant 创建租户
	CreateTenant(ctx *context.Context, req *tenant.CreateRequest) error

	// UpdateTenant 更新租户
	UpdateTenant(ctx *context.Context, req *tenant.UpdateRequest) error

	// DeleteTenant 删除租户
	DeleteTenant(ctx *context.Context, req *schema.IDRequest) error
}

// tenantService 租户服务实现
type tenantService struct {
	tenantRepo tenantrepo.Repository
	logService operate_log.OperateLogService
}

// NewTenantService 创建租户服务实例（Wire 注入）
func NewTenantService(tenantRepo tenantrepo.Repository, logService operate_log.OperateLogService) TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
		logService: logService,
	}
}

func (*tenantService) logPrefix() string {
	return "tenant-service"
}

// GetTenantByID 根据ID获取租户
func (s *tenantService) GetTenantByID(ctx *context.Context, id uint64) (*tenant.Tenant, error) {
	t, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("%s 获取租户信息失败 GetByID %d %v", s.logPrefix(), id, err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if t == nil {
		ctx.Logger.Warnf("%s 租户不存在: %d", s.logPrefix(), id)
		return nil, i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.tenant", nil)})
	}

	return t, nil
}

// ListTenants 获取租户列表
func (s *tenantService) ListTenants(ctx *context.Context, req *tenant.ListRequest) ([]*tenant.Tenant, int64, error) {
	list, total, err := s.tenantRepo.PageList(ctx, req)
	if err != nil {
		ctx.Logger.Errorf("%s 获取租户列表失败: %v", s.logPrefix(), err)
		return nil, 0, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if total == 0 {
		return []*tenant.Tenant{}, 0, nil
	}
	return list, total, nil
}

// CreateTenant 创建租户
func (s *tenantService) CreateTenant(ctx *context.Context, req *tenant.CreateRequest) error {
	// 检查租户编码是否已存在
	exists, err := s.tenantRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		ctx.Logger.Errorf("%s 检查租户编码是否存在失败: %s %v", s.logPrefix(), req.Code, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if exists {
		ctx.Logger.Warnf("%s 租户编码已存在: %s", s.logPrefix(), req.Code)
		return i18n.E(
			ctx.Context, "common.HadExist", map[string]any{"item": i18n.T(ctx.Context, "common.item.tenant", nil)})
	}

	// 设置默认状态
	status := tenant.TenantStatusEnabled
	if req.Status == 2 {
		status = tenant.TenantStatusDisabled
	}

	// 创建租户
	t := &tenant.Tenant{
		Name:         req.Name,
		Code:         req.Code,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Status:       status,
		Config:       req.Config,
	}

	err = s.tenantRepo.Create(ctx, t)
	if err != nil {
		ctx.Logger.Errorf("%s 创建租户失败: %s %v", s.logPrefix(), req.Code, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	s.logService.CreateOperateLog(ctx, i18n.T(ctx.Context, "operate.Tenant.Create", nil))

	ctx.Logger.Infof("%s 创建租户成功: %s", s.logPrefix(), req.Code)
	return nil
}

// UpdateTenant 更新租户
func (s *tenantService) UpdateTenant(ctx *context.Context, req *tenant.UpdateRequest) error {
	// 获取租户信息
	t, err := s.tenantRepo.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取租户信息失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if t == nil {
		ctx.Logger.Warnf("%s 租户不存在: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.tenant", nil)})
	}

	// 检查租户编码是否已被其他租户使用
	if req.Code != "" && req.Code != t.Code {
		exists, err := s.tenantRepo.ExistsByCode(ctx, req.Code, req.ID)
		if err != nil {
			ctx.Logger.Errorf("%s 检查租户编码是否存在失败: %s %v", s.logPrefix(), req.Code, err)
			return i18n.E(ctx.Context, "common.RepositoryErr", nil)
		}
		if exists {
			ctx.Logger.Warnf("%s 租户编码已存在: %s", s.logPrefix(), req.Code)
			return i18n.E(
				ctx.Context, "common.HadExist", map[string]any{"item": i18n.T(ctx.Context, "common.item.tenant", nil)})
		}
		t.Code = req.Code
	}

	// 更新字段
	t.Name = req.Name
	t.ContactEmail = req.ContactEmail
	t.ContactPhone = req.ContactPhone
	if req.Status == 1 {
		t.Status = tenant.TenantStatusEnabled
	} else if req.Status == 2 {
		t.Status = tenant.TenantStatusDisabled
	}
	t.Config = req.Config

	// 更新租户
	err = s.tenantRepo.Update(ctx, t)
	if err != nil {
		ctx.Logger.Errorf("%s 更新租户失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	s.logService.CreateOperateLog(ctx, i18n.T(ctx.Context, "operate.Tenant.Update", nil))

	ctx.Logger.Infof("%s 更新租户成功: %d", s.logPrefix(), req.ID)
	return nil
}

// DeleteTenant 删除租户
func (s *tenantService) DeleteTenant(ctx *context.Context, req *schema.IDRequest) error {
	// 获取租户信息
	t, err := s.tenantRepo.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取租户信息失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if t == nil {
		ctx.Logger.Warnf("%s 租户不存在: %d", s.logPrefix(), req.ID)
		return i18n.E(
			ctx.Context, "common.NotFound",
			map[string]any{"item": i18n.T(ctx.Context, "common.item.tenant", nil)})
	}

	// 删除租户
	err = s.tenantRepo.Delete(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 删除租户失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	s.logService.CreateOperateLog(ctx, i18n.T(ctx.Context, "operate.Tenant.Delete", nil))

	ctx.Logger.Infof("%s 删除租户成功: %d", s.logPrefix(), req.ID)
	return nil
}
