package context

import (
	"goadmin/pkg/logger"
	"goadmin/pkg/trace"

	"github.com/gin-gonic/gin"
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

type HandlerFunc = func(ctx *Context)

func Build(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := New(c)
		h(ctx)
	}
}

func New(c *gin.Context) *Context {
	return &Context{
		Context: c,
		Logger:  logger.Global().With(trace.GetTrace(c)),
	}
}
