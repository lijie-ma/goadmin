package role

import (
	"goadmin/internal/context"
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
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	roles, total, err := h.roleSrv.ListRoles(ctx, &pq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("GetRoleListSuccess"),
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
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	role, err := h.roleSrv.GetRoleWithPermissions(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	if role == nil {
		ctx.JSON(http.StatusNotFound, schema.Response{
			Code:    http.StatusNotFound,
			Message: ctx.Show("RoleNotFound"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("GetRoleDetailSuccess"),
		Data:    role,
	})
}

// CreateRole 创建角色
func (h *Handler) CreateRole(ctx *context.Context) {
	var req modelrole.Role
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	err := h.roleSrv.CreateRole(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("CreateRoleSuccess"),
	})
}

// UpdateRole 更新角色
func (h *Handler) UpdateRole(ctx *context.Context) {
	var req modelrole.Role
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	if req.ID == 0 {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: ctx.Show("RoleIDRequired"),
		})
		return
	}

	err := h.roleSrv.UpdateRole(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("UpdateRoleSuccess"),
	})
}

// DeleteRole 删除角色
func (h *Handler) DeleteRole(ctx *context.Context) {
	var req struct {
		ID uint64 `json:"id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	err := h.roleSrv.DeleteRole(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("DeleteRoleSuccess"),
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
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	err := h.roleSrv.AssignPermissions(ctx, req.RoleCode, req.PermissionCodes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("AssignPermissionsSuccess"),
	})
}

// GetRolePermissions 获取角色的权限列表
func (h *Handler) GetRolePermissions(ctx *context.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	permissions, err := h.roleSrv.GetRolePermissions(ctx, req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("GetRolePermissionsSuccess"),
		Data:    permissions,
	})
}

// ListActiveRoles 获取所有激活的角色
func (h *Handler) ListActiveRoles(ctx *context.Context) {
	roles, err := h.roleSrv.ListActiveRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: ctx.Show("InternalError"),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("GetActiveRoleListSuccess"),
		Data:    roles,
	})
}
