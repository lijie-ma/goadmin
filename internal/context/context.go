package context

import (
	"goadmin/pkg/logger"
	"goadmin/pkg/trace"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Context struct {
	*gin.Context
	Logger logger.Logger
}

func (c *Context) Session() Session {
	data, exists := c.Get(gin.AuthUserKey)
	if !exists {
		return nil
	}
	return data.(Session)
}

func (c *Context) ToCli() *CliContext {
	return &CliContext{
		Context: c,
		CancelFunc: func() {
			c.Abort()
		},
		Logger: c.Logger,
	}
}

func (c *Context) Show(messageID string) string {
	loc := c.MustGet("localizer").(*i18n.Localizer)
	return loc.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
}

func (c *Context) ShowWithData(messageID string, data map[string]interface{}) string {
	loc := c.MustGet("localizer").(*i18n.Localizer)
	return loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
}

type HandlerFunc = func(ctx *Context)

func Build(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			Context: c,
			Logger:  logger.Global().With(trace.GetTrace(c)),
		}
		h(ctx)
	}
}
