package middleware

import (
	"bytes"
	"goadmin/pkg/logger"
	"goadmin/pkg/trace"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 返回一个记录HTTP请求日志的中间件
func Logger() gin.HandlerFunc {
	log := logger.Global()

	return func(c *gin.Context) {
		// 请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 恢复请求体供后续中间件使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建响应体缓冲区
		responseBodyWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseBodyWriter

		// 处理请求
		c.Next()

		// 请求结束时间
		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		// 判断是否记录响应体（例如仅记录错误响应）
		var responseBody string
		if status >= 400 {
			responseBody = responseBodyWriter.body.String()
		}

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 记录请求日志
		logFields := []logger.Field{
			logger.String("method", method),
			logger.String("path", path),
			logger.String("query", query),
			logger.String("ip", clientIP),
			logger.Int("status", status),
			logger.Int("size", size),
			logger.String("latency", latency.String()),
			logger.String("user-agent", c.Request.UserAgent()),
		}

		// 根据内容类型和大小判断是否记录请求体
		contentType := c.GetHeader("Content-Type")
		if len(requestBody) > 0 && (bytes.Contains([]byte(contentType), []byte("json")) ||
			bytes.Contains([]byte(contentType), []byte("form"))) && len(requestBody) < 10240 { // 限制大小为10KB
			logFields = append(logFields, logger.String("request", string(requestBody)))
		}

		// 添加错误响应体
		if responseBody != "" && len(responseBody) < 10240 {
			logFields = append(logFields, logger.String("response", responseBody))
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			logFields = append(logFields, logger.String("errors", c.Errors.String()))
		}

		// 根据状态码选择日志级别
		logMsg := "HTTP Request"
		log = log.With(trace.GetTrace(c))
		if status >= 500 {
			log.Error(logMsg, logFields...)
		} else if status >= 400 {
			log.Warn(logMsg, logFields...)
		} else {
			log.Info(logMsg, logFields...)
		}
	}
}

// responseBodyWriter 用于捕获响应体
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应体的同时复制到缓冲区
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
