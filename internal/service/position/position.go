package position

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/position"
	"goadmin/internal/model/schema"
	positionrepo "goadmin/internal/repository/position"
	"goadmin/internal/service/operate_log"
)

// PositionService 位置服务接口
type PositionService interface {
	// GetPositionByID 根据ID获取位置
	GetPositionByID(ctx *context.Context, id uint64) (*position.Position, error)

	// ListPositions 获取位置列表
	ListPositions(ctx *context.Context, req *position.ListRequest) ([]*position.Position, int64, error)

	// CreatePosition 创建位置
	CreatePosition(ctx *context.Context, req *position.CreatePositionRequest) error

	// UpdatePosition 更新位置
	UpdatePosition(ctx *context.Context, req *position.UpdatePositionRequest) error

	// DeletePosition 删除位置
	DeletePosition(ctx *context.Context, req *schema.IDRequest) error
}

// positionService 位置服务实现
type positionService struct {
	positionRepo positionrepo.PositionRepository
	logService   operate_log.OperateLogService
}

// NewPositionService 创建位置服务实例
func NewPositionService() PositionService {
	return &positionService{
		positionRepo: positionrepo.NewPositionRepository(),
		logService:   operate_log.NewOperateLogService(),
	}
}

func (*positionService) logPrefix() string {
	return "position-service"
}

// GetPositionByID 根据ID获取位置
func (s *positionService) GetPositionByID(ctx *context.Context, id uint64) (*position.Position, error) {
	p, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("%s 获取位置信息失败 GetByID %d %v", s.logPrefix(), id, err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if p == nil {
		ctx.Logger.Warnf("%s 位置不存在: %d", s.logPrefix(), id)
		return nil, i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.position", nil)})
	}

	return p, nil
}

// ListPositions 获取位置列表
func (s *positionService) ListPositions(ctx *context.Context, req *position.ListRequest) ([]*position.Position, int64, error) {
	list, total, err := s.positionRepo.PageList(ctx, req)
	if err != nil {
		ctx.Logger.Errorf("%s 获取位置列表失败: %v", s.logPrefix(), err)
		return nil, 0, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if total == 0 {
		return []*position.Position{}, 0, nil
	}
	return list, total, nil
}

// CreatePosition 创建位置
func (s *positionService) CreatePosition(ctx *context.Context, req *position.CreatePositionRequest) error {
	// 检查位置名称是否已存在
	exists, err := s.positionRepo.IsLocationExists(ctx, req.Location)
	if err != nil {
		ctx.Logger.Errorf("%s 检查位置名称是否存在失败: %s %v", s.logPrefix(), req.Location, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if exists {
		ctx.Logger.Warnf("%s 位置名称已存在: %s", s.logPrefix(), req.Location)
		return i18n.E(
			ctx.Context, "common.HadExist", map[string]any{"item": i18n.T(ctx.Context, "common.item.position", nil)})
	}

	// 创建位置
	pos := &position.Position{
		City:       req.City,
		Location:   req.Location,
		Longitude:  req.Longitude,
		Latitude:   req.Latitude,
		CustomName: req.CustomName,
		CreatorID:  int(ctx.Session().GetID()),
		Creator:    ctx.Session().GetUsername(),
	}

	err = s.positionRepo.Create(ctx, pos)
	if err != nil {
		ctx.Logger.Errorf("%s 创建位置失败: %s %v", s.logPrefix(), req.Location, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	s.logService.CreateOperateLog(ctx, i18n.T(ctx.Context, "operate.Position.Create", nil))

	ctx.Logger.Infof("%s 创建位置成功: %s", s.logPrefix(), req.Location)
	return nil
}

// UpdatePosition 更新位置
func (s *positionService) UpdatePosition(ctx *context.Context, req *position.UpdatePositionRequest) error {
	// 获取位置信息
	pos, err := s.positionRepo.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取位置信息失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if pos == nil {
		ctx.Logger.Warnf("%s 位置不存在: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.position", nil)})
	}

	// 检查位置名称是否已被其他位置使用
	if req.Location != "" && req.Location != pos.Location {
		exists, err := s.positionRepo.IsLocationExists(ctx, req.Location, req.ID)
		if err != nil {
			ctx.Logger.Errorf("%s 检查位置名称是否存在失败: %s %v", s.logPrefix(), req.Location, err)
			return i18n.E(ctx.Context, "common.RepositoryErr", nil)
		}
		if exists {
			ctx.Logger.Warnf("%s 位置名称已存在: %s", s.logPrefix(), req.Location)
			return i18n.E(
				ctx.Context, "common.HadExist", map[string]any{"item": i18n.T(ctx.Context, "common.item.position", nil)})
		}
		pos.Location = req.Location
	}

	// 更新字段
	pos.City = req.City
	pos.Longitude = req.Longitude
	pos.Latitude = req.Latitude
	pos.CustomName = req.CustomName

	// 更新位置
	err = s.positionRepo.Update(ctx, pos)
	if err != nil {
		ctx.Logger.Errorf("%s 更新位置失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	s.logService.CreateOperateLog(ctx, i18n.T(ctx.Context, "operate.Position.Update", nil))

	ctx.Logger.Infof("%s 更新位置成功: %d", s.logPrefix(), req.ID)
	return nil
}

// DeletePosition 删除位置
func (s *positionService) DeletePosition(ctx *context.Context, req *schema.IDRequest) error {
	// 获取位置信息
	pos, err := s.positionRepo.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取位置信息失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if pos == nil {
		ctx.Logger.Warnf("%s 位置不存在: %d", s.logPrefix(), req.ID)
		return i18n.E(
			ctx.Context, "common.NotFound",
			map[string]any{"item": i18n.T(ctx.Context, "common.item.position", nil)})
	}

	// 删除位置
	err = s.positionRepo.Delete(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 删除位置失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	s.logService.CreateOperateLog(ctx, i18n.T(ctx.Context, "operate.Position.Delete", nil))

	ctx.Logger.Infof("%s 删除位置成功: %d", s.logPrefix(), req.ID)
	return nil
}
