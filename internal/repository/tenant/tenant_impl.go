package tenant

import (
	"context"
	"errors"

	"goadmin/internal/model/tenant"
	"goadmin/pkg/db"

	"gorm.io/gorm"
)

// 确保 TenantRepositoryImpl 实现了 Repository 接口
var _ Repository = (*TenantRepositoryImpl)(nil)

// TenantRepositoryImpl 实现 Repository 接口
type TenantRepositoryImpl struct {
	*db.BaseRepository[tenant.Tenant]
}

// NewTenantRepositoryImpl 创建租户仓储实例（Wire 注入）
func NewTenantRepositoryImpl(database *gorm.DB) *TenantRepositoryImpl {
	return &TenantRepositoryImpl{
		db.NewBaseRepository[tenant.Tenant](database),
	}
}

// NewTenantRepository 创建租户仓储实例（接口类型，Wire 用）
func NewTenantRepository(database *gorm.DB) Repository {
	return NewTenantRepositoryImpl(database)
}

// Deprecated: 使用 NewTenantRepository 替代
// NewTenantRepository_legacy 创建租户仓储实例（兼容旧代码，使用全局db）
func NewTenantRepository_legacy() Repository {
	return NewTenantRepositoryImpl(db.GetDB())
}

// GetByCode 根据编码获取租户
func (r *TenantRepositoryImpl) GetByCode(ctx context.Context, code string) (*tenant.Tenant, error) {
	var t tenant.Tenant
	err := r.DB().WithContext(ctx).Where("code = ?", code).First(&t).Error
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// PageList 获取租户列表（带分页和搜索）
func (r *TenantRepositoryImpl) PageList(ctx context.Context, req *tenant.ListRequest) ([]*tenant.Tenant, int64, error) {
	opts := []db.QueryOption[tenant.Tenant]{
		db.Order[tenant.Tenant](req.OrderBy),
	}

	// 如果有搜索关键词，添加搜索条件
	if req.Keyword != "" {
		opts = append(opts, db.Where[tenant.Tenant]("name LIKE ? OR code LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%"))
	}

	// 如果有状态筛选，添加状态条件
	if req.Status != nil {
		opts = append(opts, db.Where[tenant.Tenant]("status = ?", *req.Status))
	}

	return r.List(ctx, req.Page, req.PageSize, opts...)
}

// ExistsByCode 检查编码是否存在
func (r *TenantRepositoryImpl) ExistsByCode(ctx context.Context, code string, excludeID ...uint64) (bool, error) {
	opts := []db.QueryOption[tenant.Tenant]{
		db.Where[tenant.Tenant]("code = ?", code),
	}

	// 排除指定ID
	if len(excludeID) > 0 && excludeID[0] > 0 {
		opts = append(opts, db.Where[tenant.Tenant]("id != ?", excludeID[0]))
	}
	return r.Exists(ctx, opts...)
}
