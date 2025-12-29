package trace

import (
	"goadmin/pkg/logger"
	"goadmin/pkg/util"

	"github.com/gin-gonic/gin"
)

// 定义用于存储追踪信息的上下文键
const (
	TraceIDKey = "X-Request-ID"
)

// SetTraceID 设置跟踪ID
// 如果原来有跟踪ID，则使用原来的，否则生成新的
func SetTraceID(c *gin.Context) {
	traceID := GetTraceValue(c)
	if traceID == "" {
		traceID = util.GenerateUUID()
	}
	// 将跟踪ID设置到Gin上下文中
	c.Set(TraceIDKey, traceID)
	// 设置响应头中的跟踪ID
	c.Header(TraceIDKey, traceID)
}

// GetTrace 从Gin上下文获取跟踪ID
func GetTrace(c *gin.Context) logger.Field {
	return logger.String(TraceIDKey, GetTraceValue(c))
}

// GetTrace 从Gin上下文获取跟踪ID
func GetTraceValue(c *gin.Context) string {
	traceID := c.GetHeader(TraceIDKey)
	if traceID == "" {
		traceID = c.GetString(TraceIDKey)
	}
	return traceID
}

func NewTraceID() logger.Field {
	return logger.String(TraceIDKey, util.GenerateUUID())
}
