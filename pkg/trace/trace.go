package trace

import (
	"goadmin/pkg/logger"
	"goadmin/pkg/util"

	"github.com/gin-gonic/gin"
)

// 定义用于存储追踪信息的上下文键
const (
	TraceIDKey    = "X-Trace-ID"
	TraceStartKey = "X-Trace-Start"
)

// SetTraceID 设置跟踪ID
// 如果原来有跟踪ID，则使用原来的，否则生成新的
func SetTraceID(c *gin.Context) {
	traceID := c.GetHeader(TraceIDKey)
	if traceID == "" {
		traceID = util.GenerateUUID()
	}
	c.Set(TraceIDKey, traceID)
	// 设置响应头中的跟踪ID
	c.Header(TraceIDKey, traceID)
}

// GetTrace 从Gin上下文获取跟踪ID
func GetTrace(c *gin.Context) logger.Field {
	traceID, exists := c.Get(TraceIDKey)
	if !exists {
		return logger.String(TraceIDKey, "no-trace-id")
	}
	return logger.String(TraceIDKey, traceID.(string))
}
