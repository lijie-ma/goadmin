package db

import (
	"context"
)

// Model 定义了基础模型接口
type Model interface {
	TableName() string
}

// BaseRepository 基础仓储接口
type Repository[T Model] interface {
	// Create 创建记录
	Create(ctx context.Context, model *T) error

	// GetByID 根据ID获取记录
	GetByID(ctx context.Context, id uint64) (*T, error)

	// List 获取记录列表
	List(ctx context.Context, page, pageSize int) ([]*T, int64, error)

	// Update 更新记录
	Update(ctx context.Context, model *T) error

	// Delete 删除记录
	Delete(ctx context.Context, id uint64) error

	// BatchCreate 批量创建记录
	BatchCreate(ctx context.Context, models []*T) error

	// BatchDelete 批量删除记录
	BatchDelete(ctx context.Context, ids []uint64) error

	// GetByIDs 根据ID列表获取多条记录
	GetByIDs(ctx context.Context, ids []uint64) ([]*T, error)

	// Count 获取记录总数
	Count(ctx context.Context) (int64, error)
}
