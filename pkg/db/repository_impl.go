package db

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// BaseRepository 通用仓储实现
type BaseRepository[T Model] struct {
	db *gorm.DB
}

// NewBaseRepository 创建通用仓储实例
func NewBaseRepository[T Model](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

// Create 创建记录
func (r *BaseRepository[T]) Create(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

// GetByID 根据ID获取记录
func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint64) (*T, error) {
	var model T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

// applyOptions 应用查询选项
func (r *BaseRepository[T]) applyOptions(db *gorm.DB, opts ...QueryOption[T]) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

// List 获取记录列表，支持分页和自定义查询条件
func (r *BaseRepository[T]) List(ctx context.Context, page, pageSize int, opts ...QueryOption[T]) ([]*T, int64, error) {
	var (
		models []*T
		total  int64
	)

	// 构建查询
	query := r.db.WithContext(ctx).Model(new(T))
	query = r.applyOptions(query, opts...)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 没有数据直接返回
	if total == 0 {
		return models, 0, nil
	}

	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 查询数据
	err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error

	return models, total, err
}

// Find 根据条件查询记录，不分页
func (r *BaseRepository[T]) Find(ctx context.Context, opts ...QueryOption[T]) ([]*T, error) {
	var models []*T

	// 构建查询
	query := r.db.WithContext(ctx)
	query = r.applyOptions(query, opts...)

	// 查询数据
	err := query.Find(&models).Error
	return models, err
}

// Exists 根据条件判断数据是否存在
func (r *BaseRepository[T]) Exists(ctx context.Context, opts ...QueryOption[T]) (bool, error) {
	var exists bool
	sub := r.db.Model(new(T)).Select("1")
	sub = r.applyOptions(sub, opts...)
	err := r.db.WithContext(ctx).Select("EXISTS(?)", sub).Scan(&exists).Error
	return exists, err
}

// Update 更新记录
func (r *BaseRepository[T]) Update(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 删除记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(new(T)).Error
}

// BatchCreate 批量创建记录
func (r *BaseRepository[T]) BatchCreate(ctx context.Context, models []*T) error {
	if len(models) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(models).Error
}

// BatchDelete 批量删除记录
func (r *BaseRepository[T]) BatchDelete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(new(T)).Error
}

// GetByIDs 根据ID列表获取多条记录
func (r *BaseRepository[T]) GetByIDs(ctx context.Context, ids []uint64) ([]*T, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var models []*T
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&models).Error
	return models, err
}

// Count 获取记录总数，支持自定义查询条件
func (r *BaseRepository[T]) Count(ctx context.Context, opts ...QueryOption[T]) (int64, error) {
	var count int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(new(T))
	query = r.applyOptions(query, opts...)

	// 获取总数
	err := query.Count(&count).Error
	return count, err
}

// WithTx 在事务中执行操作
func (r *BaseRepository[T]) WithTx(tx *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: tx,
	}
}

// DB 获取数据库实例
func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.db
}
