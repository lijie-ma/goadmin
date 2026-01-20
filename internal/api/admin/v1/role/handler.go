package role

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	modelrole "goadmin/internal/model/role"
	"goadmin/internal/model/schema"
	rolesrv "goadmin/internal/service/role"
	"net/http"
)

// Handler 角色API处理程序
type Handler struct {
	roleSrv rolesrv.RoleService
}

// NewHandler 创建角色API处理程序
func NewHandler(roleSrv rolesrv.RoleService) *Handler {
	return &Handler{
		roleSrv: roleSrv,
	}
}

// ListRoles 获取角色列表
func (h *Handler) ListRoles(ctx *context.Context) {
	var pq schema.PageRequest
	if err := ctx.ShouldBindQuery(&pq); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	roles, total, err := h.roleSrv.ListRoles(ctx, &pq)
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
			"list":  roles,
			"total": total,
		},
	})
}

// GetRole 获取角色详情
func (h *Handler) GetRole(ctx *context.Context) {
	var req struct {
		ID uint64 `json:"id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	role, err := h.roleSrv.GetRoleWithPermissions(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if role == nil {
		ctx.JSON(http.StatusNotFound, schema.Response{
			Code:    http.StatusNotFound,
			Message: i18n.T(ctx.Context, "common.NotFound", nil),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: i18n.T(ctx.Context, "common.ActionSuccess", nil),
		Data:    role,
	})
}

// CreateRole 创建角色
func (h *Handler) CreateRole(ctx *context.Context) {
	var req modelrole.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.roleSrv.CreateRole(ctx, &req)
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

// UpdateRole 更新角色
func (h *Handler) UpdateRole(ctx *context.Context) {
	var req modelrole.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.roleSrv.UpdateRole(ctx, &req)
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

// DeleteRole 删除角色
func (h *Handler) DeleteRole(ctx *context.Context) {
	var req schema.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.roleSrv.DeleteRole(ctx, req.ID)
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

// AssignPermissions 分配权限给角色
func (h *Handler) AssignPermissions(ctx *context.Context) {
	var req struct {
		RoleCode        string   `json:"role_code" binding:"required"`
		PermissionCodes []string `json:"permission_codes" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.roleSrv.AssignPermissions(ctx, req.RoleCode, req.PermissionCodes)
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

// GetRolePermissions 获取角色的权限列表
func (h *Handler) GetRolePermissions(ctx *context.Context) {
	var req struct {
		Code string `form:"code" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	permissions, err := h.roleSrv.GetRolePermissions(ctx, req.Code)
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
		Data:    permissions,
	})
}

// ListActiveRoles 获取所有激活的角色
func (h *Handler) ListActiveRoles(ctx *context.Context) {
	roles, err := h.roleSrv.ListActiveRoles(ctx)
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
		Data:    roles,
	})
}

// ListAllPermissions 获取所有权限列表
func (h *Handler) ListAllPermissions(ctx *context.Context) {
	permissions, err := h.roleSrv.ListAllPermissions(ctx)
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
		Data:    permissions,
	})
}
