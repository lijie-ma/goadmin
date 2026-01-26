package operate_log

import (
	"context"
	"goadmin/internal/model/operate_log"
	"goadmin/pkg/db"
)

// OperateLogRepository 定义操作日志仓储接口
type OperateLogRepository interface {
	db.Repository[operate_log.OperateLog]

	// PageList 获取操作日志列表（支持多条件查询）
	PageList(ctx context.Context, req *operate_log.ListRequest) ([]*operate_log.OperateLog, int64, error)

	// CreateLog 创建操作日志
	CreateLog(ctx context.Context, log *operate_log.OperateLog) error
}
