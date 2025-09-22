package middleware

import (
	"goadmin/pkg/trace"

	"github.com/gin-gonic/gin"
)

// Trace 返回一个请求跟踪中间件，为每个请求生成唯一ID
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置响应头中的跟踪ID
		trace.SetTraceID(c)
		// 处理请求
		c.Next()
	}
}
