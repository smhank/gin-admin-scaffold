package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey 是存储在 Gin Context 中的请求ID键名
const RequestIDKey = "requestID"

// RequestIDHeader 是响应头中请求ID的字段名
const RequestIDHeader = "X-Request-ID"

// RequestIDMiddleware 为每个请求生成唯一的请求ID，并注入到上下文和响应头中
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用客户端传入的请求ID（用于链路追踪）
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			// 如果没有则生成新的 UUID
			requestID = uuid.New().String()
		}

		// 将请求ID存入 Gin Context，供后续 handler 和中间件使用
		c.Set(RequestIDKey, requestID)

		// 设置响应头，将请求ID返回给客户端
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID 从 Gin Context 中获取请求ID
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		return id.(string)
	}
	return ""
}
