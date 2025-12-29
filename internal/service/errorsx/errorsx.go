package errorsx

import (
	"errors"

	"goadmin/internal/i18n"

	"github.com/gin-gonic/gin"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid parameter")
	ErrReqired  = errors.New("required parameter")
)

type I18nError struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

func (e I18nError) Error() string { return e.Msg }

// NotFound 使用 item 的多语言 key
func New(c *gin.Context, itemKey string, ext map[string]any) I18nError {
	return I18nError{Code: itemKey, Msg: i18n.T(c, itemKey, ext)}
}
