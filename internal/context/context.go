package context

import (
	"goadmin/internal/middleware"
	"goadmin/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
	Logger logger.Logger
}

type HandlerFunc = func(ctx *Context)

func Build(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			Context: c,
			Logger: logger.Global().With(
				logger.String(middleware.TraceIDKey, middleware.GetTraceID(c))),
		}
		h(ctx)
	}
}
