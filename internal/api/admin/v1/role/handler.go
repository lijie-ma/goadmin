package role

import (
	"goadmin/internal/context"
	modelrole "goadmin/internal/model/role"
	"goadmin/internal/model/schema"
	rolesrv "goadmin/internal/service/role"
	"net/http"
	"strconv"
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
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	roles, total, err := h.roleSrv.ListRoles(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "获取角色列表成功",
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
			Message: "无效的请求参数",
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
			Message: "角色不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "获取角色详情成功",
		Data:    role,
	})
}

// CreateRole 创建角色
func (h *Handler) CreateRole(ctx *context.Context) {
	var req modelrole.Role
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
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
		Message: "创建角色成功",
	})
}

// UpdateRole 更新角色
func (h *Handler) UpdateRole(ctx *context.Context) {
	var req modelrole.Role
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
		})
		return
	}

	if req.ID == 0 {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "角色ID不能为空",
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
		Message: "更新角色成功",
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
			Message: "无效的请求参数",
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
		Message: "删除角色成功",
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
			Message: "无效的请求参数",
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
		Message: "分配权限成功",
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
			Message: "无效的请求参数",
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
		Message: "获取角色权限成功",
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
		Message: "获取激活角色列表成功",
		Data:    roles,
	})
}
