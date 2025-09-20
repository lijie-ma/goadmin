package captcha

import (
	"goadmin/internal/context"
	modelCaptcha "goadmin/internal/model/captcha"
	"goadmin/internal/service/captcha"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenerateHandler 生成验证码
func GenerateHandler(ctx *context.Context) {
	captData, err := captcha.Generate(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成验证码失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    captData,
	})
}

// CheckHandler 校验验证码
func CheckHandler(ctx *context.Context) {
	var formData modelCaptcha.CheckForm
	if err := ctx.ShouldBind(&formData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}
	err := captcha.Check(ctx, formData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "验证码错误",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}
