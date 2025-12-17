package setting

import (
	"errors"
	"goadmin/internal/context"
	"goadmin/internal/model/schema"
	"goadmin/internal/model/server"
	"goadmin/internal/service/errorsx"
	"goadmin/internal/service/setting"
	"net/http"
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
			Message: c.Show("ActionFailed"),
		})
		return
	}

	c.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: c.Show("ActionSuccess"),
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
			Message: c.Show("BadParameter"),
		})
		return
	}

	if err := h.settingSrv.SetSystemSettings(c, &settings); err != nil {
		c.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: c.ShowWithData("UpdateSettingsFailed", map[string]interface{}{"error": err.Error()}),
		})
		return
	}

	c.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: c.Show("UpdateSettingsSuccess"),
	})
}

// GetByName 根据名称获取配置
// @Summary 根据名称获取配置
// @Description 根据配置名称获取配置值
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param name query string true "配置名称"
// @Success 200 {object} schema.Response
// @Router /api/admin/v1/settings/name [get]
func (h *Handler) GetByName(ctx *context.Context) {
	var req struct {
		Name string `form:"name" json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: ctx.Show("BadParameter"),
		})
		return
	}

	var (
		rs  any
		err error
	)

	err = h.settingSrv.GetValue(ctx, req.Name, &rs)
	if err != nil {
		if errors.Is(err, errorsx.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, schema.Response{
				Code: http.StatusNotFound,
				Message: ctx.ShowWithData("NotFound", map[string]any{
					"name": ctx.Show("Config"),
				}),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: ctx.Show("ActionSuccess"),
		Data:    rs,
	})
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
			Message: ctx.Show("BadParameter"),
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
		Message: ctx.Show("ActionSuccess"),
	})
}
