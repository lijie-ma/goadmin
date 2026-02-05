package position

import (
	"context"
	"errors"
	"goadmin/internal/model/position"
	"goadmin/pkg/db"
)

// 确保PositionRepositoryImpl实现了PositionRepository接口
var _ PositionRepository = (*PositionRepositoryImpl)(nil)

// PositionRepositoryImpl 实现PositionRepository接口
type PositionRepositoryImpl struct {
	*db.BaseRepository[position.Position]
}

// NewPositionRepository 创建位置仓储实例
func NewPositionRepository() PositionRepository {
	return &PositionRepositoryImpl{
		db.NewBaseRepository[position.Position](db.GetDB()),
	}
}

// GetByID 根据ID获取位置
func (r *PositionRepositoryImpl) GetByID(ctx context.Context, id uint64) (*position.Position, error) {
	var p position.Position
	err := r.DB().WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// GetByCity 根据城市获取位置列表
func (r *PositionRepositoryImpl) GetByCity(ctx context.Context, city string) ([]*position.Position, error) {
	var positions []*position.Position
	err := r.DB().WithContext(ctx).Where("city = ?", city).Find(&positions).Error
	return positions, err
}

// PageList 获取位置列表（带分页和搜索）
func (r *PositionRepositoryImpl) PageList(ctx context.Context, req *position.ListRequest) ([]*position.Position, int64, error) {
	opts := []db.QueryOption[position.Position]{
		db.Order[position.Position](req.OrderBy),
	}

	// 如果有搜索关键词，添加搜索条件
	if req.Keyword != "" {
		opts = append(opts, db.Where[position.Position]("city LIKE ? OR location LIKE ? OR custom_name LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%"))
	}

	// 如果有城市筛选，添加城市条件
	if req.City != "" {
		opts = append(opts, db.Where[position.Position]("city = ?", req.City))
	}

	return r.List(ctx, req.Page, req.PageSize, opts...)
}

// IsLocationExists 检查位置名称是否已存在
func (r *PositionRepositoryImpl) IsLocationExists(ctx context.Context, location string, excludeID ...uint64) (bool, error) {
	var count int64
	query := r.DB().WithContext(ctx).Model(&position.Position{}).
		Where("location = ?", location)

	// 排除指定ID
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}
