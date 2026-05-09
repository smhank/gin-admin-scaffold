package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// DefaultCORSConfig 返回默认 CORS 配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// GetCORSConfig 从配置文件中读取 CORS 配置，未配置则使用默认值
func GetCORSConfig() *CORSConfig {
	cfg := DefaultCORSConfig()

	if viper.IsSet("cors.allowed_origins") {
		cfg.AllowedOrigins = viper.GetStringSlice("cors.allowed_origins")
	}
	if viper.IsSet("cors.allowed_methods") {
		cfg.AllowedMethods = viper.GetStringSlice("cors.allowed_methods")
	}
	if viper.IsSet("cors.allowed_headers") {
		cfg.AllowedHeaders = viper.GetStringSlice("cors.allowed_headers")
	}
	if viper.IsSet("cors.exposed_headers") {
		cfg.ExposedHeaders = viper.GetStringSlice("cors.exposed_headers")
	}
	if viper.IsSet("cors.allow_credentials") {
		cfg.AllowCredentials = viper.GetBool("cors.allow_credentials")
	}
	if viper.IsSet("cors.max_age") {
		cfg.MaxAge = viper.GetInt("cors.max_age")
	}

	return cfg
}

// CORSMiddleware CORS 跨域中间件，支持从配置文件读取允许的来源列表
func CORSMiddleware() gin.HandlerFunc {
	cfg := GetCORSConfig()

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 如果没有 Origin 头，则不是跨域请求，直接放行
		if origin == "" {
			c.Next()
			return
		}

		// 检查是否允许该来源
		allowedOrigin := isOriginAllowed(origin, cfg.AllowedOrigins)
		if allowedOrigin == "" {
			// 来源不被允许，如果是非简单请求（OPTIONS 预检），则拒绝
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.Next()
			return
		}

		// 设置 CORS 响应头
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))

		if len(cfg.ExposedHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposedHeaders, ", "))
		}

		if cfg.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if cfg.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", itoa(cfg.MaxAge))
		}

		// 如果是 OPTIONS 预检请求，直接返回 204
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isOriginAllowed 检查来源是否在允许列表中，返回匹配的 origin 值
func isOriginAllowed(origin string, allowedOrigins []string) string {
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return "*"
		}
		if allowed == origin {
			return origin
		}
		// 支持通配符匹配，如 https://*.example.com
		if strings.Contains(allowed, "*") {
			pattern := strings.ReplaceAll(allowed, ".", "\\.")
			pattern = strings.ReplaceAll(pattern, "*", ".*")
			if matched, _ := matchPattern(pattern, origin); matched {
				return origin
			}
		}
	}
	return ""
}

// matchPattern 简单的通配符匹配
func matchPattern(pattern, str string) (bool, error) {
	return strings.HasPrefix(str, strings.TrimSuffix(pattern, ".*")), nil
}

// itoa 整数转字符串
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
