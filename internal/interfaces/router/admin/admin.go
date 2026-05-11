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

	// 需要认证的路由分组
	authGroup := g.Group("")
	authGroup.Use(middleware.AuthMiddleware())

	// 权限管理
	authGroup.GET("/permissions", middleware.RBACMiddleware("permission:list", authSvc), h.Permission.ListPermissions)
	authGroup.POST("/permissions", middleware.RBACMiddleware("permission:create", authSvc), h.Permission.CreatePermission)
	authGroup.PUT("/permissions/:id", middleware.RBACMiddleware("permission:edit", authSvc), h.Permission.UpdatePermission)
	authGroup.DELETE("/permissions/:id", middleware.RBACMiddleware("permission:delete", authSvc), h.Permission.DeletePermission)

	// 角色管理
	authGroup.GET("/roles", middleware.RBACMiddleware("role:list", authSvc), h.Role.ListRoles)
	authGroup.POST("/roles", middleware.RBACMiddleware("role:create", authSvc), h.Role.CreateRole)
	authGroup.PUT("/roles/:id", middleware.RBACMiddleware("role:edit", authSvc), h.Role.UpdateRole)
	authGroup.DELETE("/roles/:id", middleware.RBACMiddleware("role:delete", authSvc), h.Role.DeleteRole)
	authGroup.PUT("/roles/:id/permissions", middleware.RBACMiddleware("role:edit", authSvc), h.Role.SetRolePermissions)
	authGroup.POST("/roles/assign", middleware.RBACMiddleware("role:edit", authSvc), h.Role.AssignRoleToUser)

	// 菜单管理
	authGroup.GET("/menus", middleware.RBACMiddleware("menu:list", authSvc), h.Menu.ListMenus)
	authGroup.POST("/menus", middleware.RBACMiddleware("menu:create", authSvc), h.Menu.CreateMenu)
	authGroup.PUT("/menus/:id", middleware.RBACMiddleware("menu:edit", authSvc), h.Menu.UpdateMenu)
	authGroup.DELETE("/menus/:id", middleware.RBACMiddleware("menu:delete", authSvc), h.Menu.DeleteMenu)
	authGroup.PUT("/menus/:id/roles", middleware.RBACMiddleware("menu:edit", authSvc), h.Menu.SetMenuRoles)

	// API 路径管理
	authGroup.GET("/paths", middleware.RBACMiddleware("path:list", authSvc), h.Path.ListPaths)
	authGroup.POST("/paths", middleware.RBACMiddleware("path:create", authSvc), h.Path.CreatePath)
	authGroup.PUT("/paths/:id", middleware.RBACMiddleware("path:edit", authSvc), h.Path.UpdatePath)
	authGroup.DELETE("/paths/:id", middleware.RBACMiddleware("path:delete", authSvc), h.Path.DeletePath)

	// 用户管理
	authGroup.GET("/users", middleware.RBACMiddleware("user:list", authSvc), h.User.ListUsers)
	authGroup.POST("/users", middleware.RBACMiddleware("user:create", authSvc), h.User.CreateUser)
	authGroup.PUT("/users/:id", middleware.RBACMiddleware("user:edit", authSvc), h.User.UpdateUser)
	authGroup.DELETE("/users/:id", middleware.RBACMiddleware("user:delete", authSvc), h.User.DeleteUser)
	authGroup.PUT("/users/:id/reset-password", middleware.RBACMiddleware("user:edit", authSvc), h.User.ResetPassword)
	authGroup.PUT("/users/:id/status", middleware.RBACMiddleware("user:edit", authSvc), h.User.UpdateStatus)
	authGroup.PUT("/users/:id/roles", middleware.RBACMiddleware("user:edit", authSvc), h.User.AssignRoles)

	// 操作日志
	authGroup.GET("/operation-logs", middleware.RBACMiddleware("operation-log:list", authSvc), h.OperationLog.ListOperationLogs)
	authGroup.DELETE("/operation-logs/:id", middleware.RBACMiddleware("operation-log:delete", authSvc), h.OperationLog.DeleteOperationLog)

	// 迁移记录
	authGroup.GET("/migrations", middleware.RBACMiddleware("migration:list", authSvc), h.Migration.ListMigrations)
}
