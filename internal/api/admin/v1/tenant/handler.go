package tenant

import (
	"net/http"

	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/schema"
	"goadmin/internal/model/tenant"
	tenantSrv "goadmin/internal/service/tenant"
)

type Handler struct {
	tenantSrv tenantSrv.TenantService
}

func NewHandler(tenantSrv tenantSrv.TenantService) *Handler {
	return &Handler{
		tenantSrv: tenantSrv,
	}
}

// CreateTenant 创建租户
func (h *Handler) CreateTenant(ctx *context.Context) {
	var req tenant.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.tenantSrv.CreateTenant(ctx, &req)
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

// UpdateTenant 更新租户
func (h *Handler) UpdateTenant(ctx *context.Context) {
	var req tenant.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.tenantSrv.UpdateTenant(ctx, &req)
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

// DeleteTenant 删除租户
func (h *Handler) DeleteTenant(ctx *context.Context) {
	var req schema.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.tenantSrv.DeleteTenant(ctx, &req)
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

// ListTenants 获取租户列表
func (h *Handler) ListTenants(ctx *context.Context) {
	var req tenant.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	tenants, total, err := h.tenantSrv.ListTenants(ctx, &req)
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
			"list":  tenants,
			"total": total,
		},
	})
}

// GetTenant 获取租户详情
func (h *Handler) GetTenant(ctx *context.Context) {
	var req schema.IDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	t, err := h.tenantSrv.GetTenantByID(ctx, req.ID)
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
		Data:    t,
	})
}
