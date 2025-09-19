package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 定义用于存储追踪信息的上下文键
const (
	TraceIDKey    = "X-Trace-ID"
	TraceStartKey = "X-Trace-Start"
)

// Trace 返回一个请求跟踪中间件，为每个请求生成唯一ID
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取跟踪ID，如果没有则生成新的
		traceID := c.GetHeader(TraceIDKey)
		if traceID == "" {
			traceID = generateTraceID()
		}

		// 记录开始时间
		startTime := time.Now()

		// 将跟踪ID和开始时间存入上下文
		c.Set(TraceIDKey, traceID)
		c.Set(TraceStartKey, startTime)

		// 设置响应头中的跟踪ID
		c.Header(TraceIDKey, traceID)

		// 处理请求
		c.Next()
	}
}

// GetTraceID 从Gin上下文获取跟踪ID
func GetTraceID(c *gin.Context) string {
	traceID, exists := c.Get(TraceIDKey)
	if !exists {
		return "no-trace-id"
	}
	return traceID.(string)
}

// 生成唯一的跟踪ID
func generateTraceID() string {
	v, _ := uuid.NewV7()
	return v.String()
}
