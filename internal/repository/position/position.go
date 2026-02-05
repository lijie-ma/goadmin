package position

import (
	"context"
	"goadmin/internal/model/position"
	"goadmin/pkg/db"
)

// PositionRepository 定义位置仓储接口
type PositionRepository interface {
	db.Repository[position.Position]

	// GetByCity 根据城市获取位置列表
	GetByCity(ctx context.Context, city string) ([]*position.Position, error)

	// PageList 获取位置列表（带分页和搜索）
	PageList(ctx context.Context, req *position.ListRequest) ([]*position.Position, int64, error)

	// IsLocationExists 检查位置名称是否已存在
	IsLocationExists(ctx context.Context, location string, excludeID ...uint64) (bool, error)
}
