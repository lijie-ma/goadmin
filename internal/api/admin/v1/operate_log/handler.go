package operate_log

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	modeloperatelog "goadmin/internal/model/operate_log"
	"goadmin/internal/model/schema"
	operatelog "goadmin/internal/service/operate_log"
	"net/http"
)

type Handler struct {
	logSrv operatelog.OperateLogService
}

func NewHandler(logSrv operatelog.OperateLogService) *Handler {
	return &Handler{
		logSrv: logSrv,
	}
}

// ListOperateLogs 获取操作日志列表
func (h *Handler) ListOperateLogs(ctx *context.Context) {
	var req modeloperatelog.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	logs, total, err := h.logSrv.ListOperateLogs(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: i18n.T(ctx.Context, "common.ActionSuccess", nil),
		Data: map[string]interface{}{
			"list":  logs,
			"total": total,
		},
	})
}
