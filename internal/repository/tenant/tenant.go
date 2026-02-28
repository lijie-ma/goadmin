package tenant

import (
	"context"

	"goadmin/internal/model/tenant"
	"goadmin/pkg/db"
)

// Repository 租户仓储接口
type Repository interface {
	db.Repository[tenant.Tenant]

	// GetByCode 根据编码获取租户
	GetByCode(ctx context.Context, code string) (*tenant.Tenant, error)

	// PageList 获取租户列表（带分页和搜索）
	PageList(ctx context.Context, req *tenant.ListRequest) ([]*tenant.Tenant, int64, error)

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, code string, excludeID ...uint64) (bool, error)
}
