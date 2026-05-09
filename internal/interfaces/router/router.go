package router

import (
	"gin-admin-base/internal/application"
	"gin-admin-base/internal/interfaces/handler"
	"gin-admin-base/internal/interfaces/middleware"
	"gin-admin-base/internal/interfaces/router/admin"
	"gin-admin-base/internal/interfaces/router/mini"

	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
	engine   *gin.Engine
	handlers *handler.Handlers
	authSvc  *application.AuthService
}

// New 创建路由管理器
func New(engine *gin.Engine, h *handler.Handlers, authSvc *application.AuthService) *Router {
	return &Router{
		engine:   engine,
		handlers: h,
		authSvc:  authSvc,
	}
}

// Register 注册所有路由
func (r *Router) Register() {
	// 管理后台端 - /api/admin，带操作日志中间件和 RBAC 权限校验
	adminGroup := r.engine.Group("/api/admin")
	adminGroup.Use(middleware.OperationLogMiddleware())
	adminGroup.Use(middleware.RateLimitMiddleware())
	admin.Register(adminGroup, r.handlers, r.authSvc)

	// 小程序端 - /api/mini
	miniGroup := r.engine.Group("/api/mini")
	miniGroup.Use(middleware.RateLimitMiddleware())
	mini.Register(miniGroup, r.handlers)
}
