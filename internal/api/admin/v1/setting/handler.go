package setting

import (
	"goadmin/internal/context"
	"goadmin/internal/model/schema"
	"goadmin/internal/model/server"
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

// GetCaptchaSwitch 获取验证码开关配置
func (h *Handler) GetCaptchaSwitch(ctx *context.Context) {
	config, err := setting.GetCaptchaSwitch(ctx, h.settingSrv)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "获取验证码开关配置成功",
		Data:    config,
	})
}

// SetCaptchaSwitch 设置验证码开关配置
func (h *Handler) SetCaptchaSwitch(ctx *context.Context) {
	var config server.CaptchaSwitchConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
		})
		return
	}

	err := setting.SetCaptchaSwitch(ctx, h.settingSrv, &config)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "设置验证码开关配置成功",
	})
}

// GetByName 根据名称获取配置
func (h *Handler) GetByName(ctx *context.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
		})
		return
	}

	setting, err := h.settingSrv.GetByName(ctx.Request.Context(), req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if setting == nil {
		ctx.JSON(http.StatusNotFound, schema.Response{
			Code:    http.StatusNotFound,
			Message: "配置不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "获取配置成功",
		Data:    setting,
	})
}

// SetByName 设置配置值
func (h *Handler) SetByName(ctx *context.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
		})
		return
	}

	err := h.settingSrv.SetByName(ctx.Request.Context(), req.Name, req.Value)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "设置配置成功",
	})
}

// BatchGetValues 批量获取配置值
func (h *Handler) BatchGetValues(ctx *context.Context) {
	var req struct {
		Names []string `json:"names" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
		})
		return
	}

	values, err := h.settingSrv.BatchGetValues(ctx.Request.Context(), req.Names)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "批量获取配置成功",
		Data:    values,
	})
}
