package middleware

import (
	"gin-admin-base/internal/application"
	"gin-admin-base/internal/interfaces/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware 权限校验中间件
// permission: 所需权限标识，如 "user:list"
// authSvc: 认证服务，nil 时跳过校验（mock 模式）
func RBACMiddleware(permission string, authSvc *application.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// mock 模式或无认证服务时跳过权限校验
		if authSvc == nil {
			c.Next()
			return
		}

		// 从上下文中获取用户 ID
		userIDStr, exists := c.Get("userID")
		if !exists {
			response.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		// 转换 userID 为 uint
		var userID uint
		switch v := userIDStr.(type) {
		case string:
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				response.Unauthorized(c, "无效的用户身份")
				c.Abort()
				return
			}
			userID = uint(id)
		case uint:
			userID = v
		case float64:
			userID = uint(v)
		default:
			response.Unauthorized(c, "无效的用户身份")
			c.Abort()
			return
		}

		// 校验权限
		hasPermission, err := authSvc.CheckPermission(userID, permission)
		if err != nil || !hasPermission {
			response.Forbidden(c, "权限不足："+permission)
			c.Abort()
			return
		}

		c.Next()
	}
}
