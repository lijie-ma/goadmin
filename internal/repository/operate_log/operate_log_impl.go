package operate_log

import (
	"context"
	"goadmin/internal/model/operate_log"
	"goadmin/pkg/db"
)

// 确保OperateLogRepositoryImpl实现了OperateLogRepository接口
var _ OperateLogRepository = (*OperateLogRepositoryImpl)(nil)

// OperateLogRepositoryImpl 实现OperateLogRepository接口
type OperateLogRepositoryImpl struct {
	*db.BaseRepository[operate_log.OperateLog]
}

// NewOperateLogRepository 创建操作日志仓储实例
func NewOperateLogRepository() OperateLogRepository {
	return &OperateLogRepositoryImpl{
		db.NewBaseRepository[operate_log.OperateLog](db.GetDB()),
	}
}

// PageList 获取操作日志列表（支持多条件查询）
func (r *OperateLogRepositoryImpl) PageList(ctx context.Context, req *operate_log.ListRequest) ([]*operate_log.OperateLog, int64, error) {
	opts := []db.QueryOption[operate_log.OperateLog]{
		db.Order[operate_log.OperateLog](req.OrderBy),
	}

	// 如果有用户名，添加查询条件
	if req.Username != "" {
		opts = append(opts, db.Where[operate_log.OperateLog]("username LIKE ?", "%"+req.Username+"%"))
	}

	// 如果有内容，添加查询条件
	if req.Content != "" {
		opts = append(opts, db.Where[operate_log.OperateLog]("content LIKE ?", "%"+req.Content+"%"))
	}

	// 如果有IP地址，添加查询条件
	if req.IP != "" {
		opts = append(opts, db.Where[operate_log.OperateLog]("ip LIKE ?", "%"+req.IP+"%"))
	}

	// 如果有开始时间，添加查询条件
	if req.StartTime != "" {
		opts = append(opts, db.Where[operate_log.OperateLog]("ctime >= ?", req.StartTime))
	}

	// 如果有结束时间，添加查询条件
	if req.EndTime != "" {
		opts = append(opts, db.Where[operate_log.OperateLog]("ctime <= ?", req.EndTime))
	}

	return r.List(ctx, req.Page, req.PageSize, opts...)
}

// CreateLog 创建操作日志
func (r *OperateLogRepositoryImpl) CreateLog(ctx context.Context, log *operate_log.OperateLog) error {
	return r.Create(ctx, log)
}
