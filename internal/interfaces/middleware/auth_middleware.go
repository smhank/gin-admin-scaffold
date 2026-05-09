package middleware

import (
	"gin-admin-base/internal/interfaces/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟 JWT 校验，实际应解析 Token 获取 UserID
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			response.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
