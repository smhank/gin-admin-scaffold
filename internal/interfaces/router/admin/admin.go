package admin

import (
	"gin-admin-base/internal/application"
	"gin-admin-base/internal/interfaces/handler"
	"gin-admin-base/internal/interfaces/middleware"

	"github.com/gin-gonic/gin"
)

// Register 注册管理后台端所有路由
func Register(g *gin.RouterGroup, h *handler.Handlers, authSvc *application.AuthService) {
	// 认证（无需权限校验）
	g.POST("/login", h.Auth.Login)

	// 权限管理
	g.GET("/permissions", middleware.RBACMiddleware("permission:list", authSvc), h.Permission.ListPermissions)
	g.POST("/permissions", middleware.RBACMiddleware("permission:create", authSvc), h.Permission.CreatePermission)
	g.PUT("/permissions/:id", middleware.RBACMiddleware("permission:edit", authSvc), h.Permission.UpdatePermission)
	g.DELETE("/permissions/:id", middleware.RBACMiddleware("permission:delete", authSvc), h.Permission.DeletePermission)

	// 角色管理
	g.GET("/roles", middleware.RBACMiddleware("role:list", authSvc), h.Role.ListRoles)
	g.POST("/roles", middleware.RBACMiddleware("role:create", authSvc), h.Role.CreateRole)
	g.PUT("/roles/:id", middleware.RBACMiddleware("role:edit", authSvc), h.Role.UpdateRole)
	g.DELETE("/roles/:id", middleware.RBACMiddleware("role:delete", authSvc), h.Role.DeleteRole)
	g.PUT("/roles/:id/permissions", middleware.RBACMiddleware("role:edit", authSvc), h.Role.SetRolePermissions)
	g.POST("/roles/assign", middleware.RBACMiddleware("role:edit", authSvc), h.Role.AssignRoleToUser)

	// 菜单管理
	g.GET("/menus", middleware.RBACMiddleware("menu:list", authSvc), h.Menu.ListMenus)
	g.POST("/menus", middleware.RBACMiddleware("menu:create", authSvc), h.Menu.CreateMenu)
	g.PUT("/menus/:id", middleware.RBACMiddleware("menu:edit", authSvc), h.Menu.UpdateMenu)
	g.DELETE("/menus/:id", middleware.RBACMiddleware("menu:delete", authSvc), h.Menu.DeleteMenu)
	g.PUT("/menus/:id/roles", middleware.RBACMiddleware("menu:edit", authSvc), h.Menu.SetMenuRoles)

	// API 路径管理
	g.GET("/paths", middleware.RBACMiddleware("path:list", authSvc), h.Path.ListPaths)
	g.POST("/paths", middleware.RBACMiddleware("path:create", authSvc), h.Path.CreatePath)
	g.PUT("/paths/:id", middleware.RBACMiddleware("path:edit", authSvc), h.Path.UpdatePath)
	g.DELETE("/paths/:id", middleware.RBACMiddleware("path:delete", authSvc), h.Path.DeletePath)

	// 用户管理
	g.GET("/users", middleware.RBACMiddleware("user:list", authSvc), h.User.ListUsers)
	g.POST("/users", middleware.RBACMiddleware("user:create", authSvc), h.User.CreateUser)
	g.PUT("/users/:id", middleware.RBACMiddleware("user:edit", authSvc), h.User.UpdateUser)
	g.DELETE("/users/:id", middleware.RBACMiddleware("user:delete", authSvc), h.User.DeleteUser)
	g.PUT("/users/:id/reset-password", middleware.RBACMiddleware("user:edit", authSvc), h.User.ResetPassword)
	g.PUT("/users/:id/status", middleware.RBACMiddleware("user:edit", authSvc), h.User.UpdateStatus)
	g.PUT("/users/:id/roles", middleware.RBACMiddleware("user:edit", authSvc), h.User.AssignRoles)

	// 操作日志
	g.GET("/operation-logs", middleware.RBACMiddleware("operation-log:list", authSvc), h.OperationLog.ListOperationLogs)
	g.DELETE("/operation-logs/:id", middleware.RBACMiddleware("operation-log:delete", authSvc), h.OperationLog.DeleteOperationLog)

	// 迁移记录
	g.GET("/migrations", middleware.RBACMiddleware("migration:list", authSvc), h.Migration.ListMigrations)
}
