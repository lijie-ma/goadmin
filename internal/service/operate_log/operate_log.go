package operate_log

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	modeloperatelog "goadmin/internal/model/operate_log"
	operatelogrepo "goadmin/internal/repository/operate_log"
)

// OperateLogService 操作日志服务接口
type OperateLogService interface {
	// ListOperateLogs 获取操作日志列表
	ListOperateLogs(ctx *context.Context, req *modeloperatelog.ListRequest) ([]*modeloperatelog.OperateLog, int64, error)

	// CreateOperateLog 创建操作日志
	//
	// @param operator  操作人
	CreateOperateLog(ctx *context.Context, content string, operator ...string) error
}

// operateLogService 操作日志服务实现
type operateLogService struct {
	logRepo operatelogrepo.OperateLogRepository
}

// NewOperateLogService 创建操作日志服务实例
func NewOperateLogService() OperateLogService {
	return &operateLogService{
		logRepo: operatelogrepo.NewOperateLogRepository(),
	}
}

func (*operateLogService) logPrefix() string {
	return "operate-log-service"
}

// ListOperateLogs 获取操作日志列表
func (s *operateLogService) ListOperateLogs(ctx *context.Context, req *modeloperatelog.ListRequest) ([]*modeloperatelog.OperateLog, int64, error) {
	list, total, err := s.logRepo.PageList(ctx, req)
	if err != nil {
		ctx.Logger.Errorf("%s 获取操作日志列表失败: %v", s.logPrefix(), err)
		return nil, 0, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if total == 0 {
		return []*modeloperatelog.OperateLog{}, 0, nil
	}
	return list, total, nil
}

// CreateOperateLog 创建操作日志
func (s *operateLogService) CreateOperateLog(ctx *context.Context, content string, operator ...string) error {
	// 获取客户端IP
	clientIP := ctx.ClientIP()

	// 获取用户名
	username := ""
	if len(operator) > 0 {
		username = operator[0]
	} else {
		session := ctx.Session()
		if session != nil {
			username = session.GetUsername()
		}
	}

	// 创建操作日志
	log := &modeloperatelog.OperateLog{
		Content:  content,
		Username: username,
		IP:       clientIP,
	}

	err := s.logRepo.CreateLog(ctx, log)
	if err != nil {
		ctx.Logger.Errorf("%s 创建操作日志失败: %v", s.logPrefix(), err)
		return err
	}

	ctx.Logger.Infof("%s 创建操作日志成功: %s", s.logPrefix(), content)
	return nil
}
