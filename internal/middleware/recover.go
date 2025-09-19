package middleware

import (
	"fmt"
	"goadmin/pkg/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Recovery 返回一个恢复中间件，用于捕获任何panic并恢复
func Recovery() gin.HandlerFunc {
	log := logger.Global()
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查连接是否已断开
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 获取堆栈信息
				stack := string(debug.Stack())

				if brokenPipe {
					log.Error("broken pipe",
						logger.String("time", time.Now().Format(time.RFC3339)),
						logger.Any("error", err),
						logger.String("request", string(httpRequest)),
					)

					// 如果客户端已经断开连接，我们无法写入响应
					c.Error(fmt.Errorf("%v", err))
					c.Abort()
					return
				}

				// 记录错误详情
				log.Error("[Recovery from panic]",
					logger.String("time", time.Now().Format(time.RFC3339)),
					logger.Any("error", err),
					logger.String("request", string(httpRequest)),
					logger.String("stack", stack),
				)

				// 返回500错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprintf("Internal Server Error: %v", err),
				})
			}
		}()

		c.Next()
	}
}
