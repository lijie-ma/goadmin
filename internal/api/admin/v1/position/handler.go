package position

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/position"
	"goadmin/internal/model/schema"
	positionSrv "goadmin/internal/service/position"
	"net/http"
)

type Handler struct {
	positionSrv positionSrv.PositionService
}

func NewHandler(positionSrv positionSrv.PositionService) *Handler {
	return &Handler{
		positionSrv: positionSrv,
	}
}

// CreatePosition 创建位置
func (h *Handler) CreatePosition(ctx *context.Context) {
	var req position.CreatePositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.positionSrv.CreatePosition(ctx, &req)
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
	})
}

// UpdatePosition 更新位置
func (h *Handler) UpdatePosition(ctx *context.Context) {
	var req position.UpdatePositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.positionSrv.UpdatePosition(ctx, &req)
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
	})
}

// DeletePosition 删除位置
func (h *Handler) DeletePosition(ctx *context.Context) {
	var req schema.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.positionSrv.DeletePosition(ctx, &req)
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
	})
}

// ListPositions 获取位置列表
func (h *Handler) ListPositions(ctx *context.Context) {
	var req position.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	positions, total, err := h.positionSrv.ListPositions(ctx, &req)
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
			"list":  positions,
			"total": total,
		},
	})
}

// GetPosition 获取位置详情
func (h *Handler) GetPosition(ctx *context.Context) {
	var req schema.IDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	position, err := h.positionSrv.GetPositionByID(ctx, req.ID)
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
		Data:    position,
	})
}
