package db

import (
	"context"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

// Model 定义了基础模型接口
type Model interface {
	TableName() string
}

// QueryOption 查询选项
type QueryOption[T Model] func(*gorm.DB) *gorm.DB

// Where 条件查询选项
func Where[T Model](query any, args ...any) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

// Order 排序查询选项
func Order[T Model](value string) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	}
}

// Preload 预加载关联查询选项
func Preload[T Model](query string, args ...any) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

// Joins 关联查询选项
func Joins[T Model](query string, args ...any) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}
}

// Select 字段选择查询选项
func Select[T Model](query any, args ...any) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(query, args...)
	}
}

// BaseRepository 基础仓储接口
type Repository[T Model] interface {
	// Create 创建记录
	Create(ctx context.Context, model *T) error

	// GetByID 根据ID获取记录
	GetByID(ctx context.Context, id uint64) (*T, error)

	// List 获取记录列表，支持分页和自定义查询条件
	List(ctx context.Context, page, pageSize int, opts ...QueryOption[T]) ([]*T, int64, error)

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

	// Count 获取记录总数，支持自定义查询条件
	Count(ctx context.Context, opts ...QueryOption[T]) (int64, error)

	// Find 根据条件查询记录，不分页
	Find(ctx context.Context, opts ...QueryOption[T]) ([]*T, error)
}
