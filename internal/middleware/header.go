package middleware

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// HeaderConfig 头部中间件配置
type HeaderConfig struct {
	// 启用CORS
	EnableCORS bool
	// 允许的域名列表
	AllowOrigins []string
	// 允许的HTTP方法
	AllowMethods []string
	// 允许的HTTP头
	AllowHeaders []string
	// 暴露的HTTP头
	ExposeHeaders []string
	// 允许凭证
	AllowCredentials bool
	// 预检请求缓存时间（秒）
	MaxAge int
	// 服务器名称（用于X-Powered-By头）
	ServerName string
	// 是否禁用"X-Powered-By"头
	DisablePoweredBy bool
	// 是否启用安全头部
	EnableSecureHeader bool
	// 额外的自定义头部
	ExtraHeaders map[string]string
}

// DefaultHeaderConfig 默认头部配置
func DefaultHeaderConfig() HeaderConfig {
	return HeaderConfig{
		EnableCORS:         true,
		AllowOrigins:       []string{"*"},
		AllowMethods:       []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:       []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-Trace-ID"},
		ExposeHeaders:      []string{"Content-Length", "X-Trace-ID"},
		AllowCredentials:   true,
		MaxAge:             86400,
		ServerName:         "GoAdmin",
		DisablePoweredBy:   false,
		EnableSecureHeader: true,
		ExtraHeaders:       make(map[string]string),
	}
}

// Header 返回一个处理HTTP头的中间件
func Header(config ...HeaderConfig) gin.HandlerFunc {
	// 使用默认配置或传入的配置
	var cfg HeaderConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultHeaderConfig()
	}

	return func(c *gin.Context) {
		// 处理CORS
		if cfg.EnableCORS {
			setCORSHeaders(c, cfg)
		}

		// 如果是预检请求，直接返回200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// 添加服务器信息
		if !cfg.DisablePoweredBy {
			c.Header("X-Powered-By", cfg.ServerName)
		}

		// 添加安全头部
		if cfg.EnableSecureHeader {
			setSecurityHeaders(c)
		}

		// 添加额外的自定义头部
		for key, value := range cfg.ExtraHeaders {
			c.Header(key, value)
		}

		// 继续处理请求
		c.Next()
	}
}

// 设置CORS头部
func setCORSHeaders(c *gin.Context, cfg HeaderConfig) {
	origin := c.Request.Header.Get("Origin")

	// 如果允许所有源或者请求源在允许列表中
	if contains(cfg.AllowOrigins, "*") || (origin != "" && contains(cfg.AllowOrigins, origin)) {
		c.Header("Access-Control-Allow-Origin", origin)
	} else if len(cfg.AllowOrigins) > 0 {
		c.Header("Access-Control-Allow-Origin", cfg.AllowOrigins[0])
	}

	// 允许凭证
	if cfg.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	// 设置允许的HTTP方法
	if len(cfg.AllowMethods) > 0 {
		c.Header("Access-Control-Allow-Methods", join(cfg.AllowMethods))
	}

	// 设置允许的HTTP头
	if len(cfg.AllowHeaders) > 0 {
		c.Header("Access-Control-Allow-Headers", join(cfg.AllowHeaders))
	}

	// 设置暴露的HTTP头
	if len(cfg.ExposeHeaders) > 0 {
		c.Header("Access-Control-Expose-Headers", join(cfg.ExposeHeaders))
	}

	// 设置预检请求缓存时间
	if cfg.MaxAge > 0 {
		c.Header("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
	}
}

// 设置安全相关的头部
func setSecurityHeaders(c *gin.Context) {
	// 阻止浏览器探测MIME类型
	c.Header("X-Content-Type-Options", "nosniff")

	// 启用XSS过滤
	c.Header("X-XSS-Protection", "1; mode=block")

	// 阻止在iframe中加载
	c.Header("X-Frame-Options", "SAMEORIGIN")

	// 设置引用策略
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

	// 设置内容安全策略（实际项目中应根据需要配置）
	// c.Header("Content-Security-Policy", "default-src 'self'")
}

// 检查字符串切片是否包含指定值
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// 将字符串切片连接为逗号分隔的字符串
func join(slice []string) string {
	return strings.Join(slice, ", ")
}
