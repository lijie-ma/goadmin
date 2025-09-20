package captcha

import (
	"goadmin/internal/context"
	modelCaptcha "goadmin/internal/model/captcha"
	"goadmin/internal/model/schema"
	"goadmin/internal/service/captcha"
	"goadmin/internal/service/token"
	"net/http"
	"time"
)

// GenerateHandler 生成验证码
func GenerateHandler(ctx *context.Context) {
	captData, err := captcha.Generate(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: "生成验证码失败",
		})
		return
	}
	ctx.JSON(http.StatusOK,
		schema.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data:    captData,
		})
}

// CheckHandler 校验验证码
func CheckHandler(ctx *context.Context) {
	var formData modelCaptcha.CheckForm
	if err := ctx.ShouldBind(&formData); err != nil {
		ctx.JSON(http.StatusBadRequest,
			schema.Response{
				Code:    http.StatusBadRequest,
				Message: "参数错误",
			})
		return
	}
	err := captcha.Check(ctx, formData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			schema.Response{
				Code:    http.StatusBadRequest,
				Message: "验证码错误",
			})
		return
	}

	tok, err := token.NewTokenService().GenerateToken(ctx, time.Second*15)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			schema.Response{
				Code:    http.StatusInternalServerError,
				Message: "生成token失败",
			})

		return
	}

	ctx.JSON(http.StatusOK,
		schema.Response{
			Code:    http.StatusOK,
			Message: "StatusOK",
			Data: map[string]any{
				"token": tok,
			},
		})
}
