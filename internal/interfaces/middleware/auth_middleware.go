package middleware

import (
	"gin-admin-base/internal/infras/config"
	"gin-admin-base/internal/interfaces/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 头获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		// 去除 Bearer 前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// 没有 Bearer 前缀，直接使用原值（兼容旧格式）
		}

		// 解析 JWT Token
		jwtCfg := config.GetJWTConfig()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtCfg.Secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "无效的访问令牌")
			c.Abort()
			return
		}

		// 从 claims 中提取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "无效的访问令牌")
			c.Abort()
			return
		}

		username, _ := claims["username"].(string)
		c.Set("username", username)
		c.Set("userID", username)
		c.Next()
	}
}
