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

// List 获取记录列表
func (r *BaseRepository[T]) List(ctx context.Context, page, pageSize int) ([]*T, int64, error) {
	var (
		models []*T
		total  int64
	)

	// 获取总数
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
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
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&models).Error

	return models, total, err
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

// Count 获取记录总数
func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error
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
