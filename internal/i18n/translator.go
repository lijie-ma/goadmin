package i18n

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// T 根据上下文翻译 MessageID
func T(c *gin.Context, messageID string, data map[string]any) string {
	localizer, ok := c.Get("localizer")
	if !ok {
		return messageID
	}

	l := localizer.(*i18n.Localizer)
	msg, err := l.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return msg
}

// E 根据上下文翻译 MessageID
func E(c *gin.Context, messageID string, data map[string]any) error {
	return errors.New(T(c, messageID, data))
}

// Translate 手动翻译（用于非 HTTP 场景：cron、worker 等）
func Translate(lang, messageID string, data map[string]any) string {
	localizer := i18n.NewLocalizer(Bundle, lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return msg
}
