package user

import (
	"goadmin/internal/context"
	"goadmin/internal/model/schema"
	modeluser "goadmin/internal/model/user"
	"goadmin/internal/repository/user"
	"goadmin/internal/service/token"
	userSrv "goadmin/internal/service/user"
	"net/http"
)

type Handler struct {
	userSrv  userSrv.UserService
	userRepo user.UserRepository
	tokenSrv *token.TokenService
}

func NewHandler(userSrv userSrv.UserService, userRepo user.UserRepository, tokenSrv *token.TokenService) *Handler {
	return &Handler{
		userSrv:  userSrv,
		userRepo: userRepo,
		tokenSrv: tokenSrv,
	}
}

// Login 用户登录
func (h *Handler) Login(ctx *context.Context) {
	var req modeluser.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
		})
		return
	}

	resp, err := h.userSrv.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "登录成功",
		Data:    resp,
	})
}

// Logout 用户退出
func (h *Handler) Logout(ctx *context.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "未提供token",
		})
		return
	}

	// 从Bearer token中提取token值
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.tokenSrv.DeleteToken(ctx, token); err != nil {
		ctx.Logger.Errorf("删除token失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "退出成功",
	})
}

// ChangePassword 修改密码
func (h *Handler) ChangePassword(ctx *context.Context) {
	var req modeluser.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
		})
		return
	}

	err := h.userSrv.ChangePassword(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: "密码修改成功",
	})
}
