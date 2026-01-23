package upload

import (
	"fmt"
	"goadmin/config"
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/schema"
	"goadmin/pkg/util"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uploadCfg *config.UploadConfig
}

func NewHandler() *Handler {
	return &Handler{
		uploadCfg: &config.Get().Upload,
	}
}

// UploadFile 上传单个文件
func (h *Handler) UploadFile(ctx *context.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "upload.fileNotFound", nil),
		})
		return
	}

	// 检查文件大小
	if file.Size > h.uploadCfg.MaxSize {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "upload.fileTooLarge", nil),
		})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !slices.Contains(h.uploadCfg.AllowedTypes, ext) {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: i18n.T(ctx.Context, "upload.fileTypeNotAllowed", nil),
		})
		return
	}

	// 生成新的文件名
	uuidStr, err := util.UUIDV7Str()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, schema.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	newFileName := uuidStr + ext

	// 创建上传目录（按日期分类）
	datePath := time.Now().Format("2006/01/02")
	fullUploadPath := filepath.Join(h.uploadCfg.Path, datePath)
	if err := os.MkdirAll(fullUploadPath, 0755); err != nil {
		ctx.Logger.Errorf("创建上传目录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: i18n.T(ctx.Context, "common.SystemError", nil),
		})
		return
	}

	// 保存文件
	filePath := filepath.Join(fullUploadPath, newFileName)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.Logger.Errorf("保存文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, schema.Response{
			Code:    http.StatusInternalServerError,
			Message: i18n.T(ctx.Context, "upload.saveFailed", nil),
		})
		return
	}

	// 返回文件信息
	fileURL := fmt.Sprintf("%s/%s", fullUploadPath, newFileName)
	ctx.JSON(http.StatusOK, schema.Response{
		Code:    http.StatusOK,
		Message: i18n.T(ctx.Context, "upload.success", nil),
		Data: gin.H{
			"filename":     newFileName,
			"original":     file.Filename,
			"size":         file.Size,
			"url":          fileURL,
			"path":         filePath,
			"content_type": file.Header.Get("Content-Type"),
		},
	})
}
