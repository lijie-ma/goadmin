package setting

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/schema"
	"goadmin/internal/model/server"
	"goadmin/internal/service/setting"
	"net/http"
	"strings"
)

// Handler 系统设置API处理程序
type Handler struct {
	settingSrv setting.ServerSettingService
}

// NewHandler 创建系统设置API处理程序
func NewHandler(settingSrv setting.ServerSettingService) *Handler {
	return &Handler{
		settingSrv: settingSrv,
	}
}

// GetSettings 获取系统设置
// @Summary 获取系统设置
// @Description 获取所有系统设置信息
// @Tags 系统设置
// @Accept  json
// @Produce json
// @Success 200 {object} SystemSettings
// @Router /api/admin/v1/settings [get]
func (h *Handler) GetSettings(c *context.Context) {
	settings, err := h.settingSrv.GetSystemSettings(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: i18n.T(c.Context, "common.ActionSuccess", nil),
		Data:    settings,
	})
}

// UpdateSettings 更新系统设置
// @Summary 更新系统设置
// @Description 更新系统设置信息
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param settings body SystemSettings true "系统设置"
// @Success 200 {object} schema.Response
// @Router /api/admin/v1/settings [put]
func (h *Handler) UpdateSettings(c *context.Context) {
	var settings server.SystemSettingsRequest
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(c.Context, "common.BadParameter", nil),
		})
		return
	}

	if err := h.settingSrv.SetSystemSettings(c, &settings); err != nil {
		c.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: i18n.T(c.Context, "common.ActionSuccess", nil),
	})
}

// GetByNames 根据名称批量获取配置
// @Summary 根据名称批量获取配置
// @Description 根据配置名称批量获取配置值
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param names query string true "配置名称，多个用逗号分隔"
// @Success 200 {object} schema.Response
// @Router /api/admin/v1/settings/get [get]
func (h *Handler) GetByNames(ctx *context.Context) {
	var req struct {
		Names string `form:"names" json:"names" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	// 解析名称列表
	names := parseNames(req.Names)
	if len(names) == 0 {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	// 批量获取配置
	result, err := h.settingSrv.GetValues(ctx, names)
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
		Data:    result,
	})
}

// parseNames 解析名称字符串，支持逗号分隔
func parseNames(namesStr string) []string {
	var names []string
	// 简单的逗号分隔解析
	for _, name := range strings.Split(namesStr, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

// SetByName 设置配置值
// @Summary 设置配置值
// @Description 根据配置名称设置配置值
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param setting body object true "配置信息"
// @Success 200 {object} schema.Response
// @Router /api/admin/v1/settings/name [post]
func (h *Handler) SetByName(ctx *context.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Value any    `json:"value" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "common.BadParameter", nil),
		})
		return
	}

	err := h.settingSrv.SetByName(ctx, req.Name, req.Value)
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
