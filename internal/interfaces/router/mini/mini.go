package mini

import (
	"gin-admin-base/internal/interfaces/handler"

	"github.com/gin-gonic/gin"
)

// Register 注册小程序端所有路由
func Register(g *gin.RouterGroup, h *handler.Handlers) {
	// 认证
	g.POST("/login", h.Auth.Login)
}
