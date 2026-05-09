package handler

import (
	"gin-admin-base/internal/application"
	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/infras/config"
	"gin-admin-base/internal/infras/global"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Handlers 聚合所有 Handler
type Handlers struct {
	Auth         *AuthHandler
	Permission   *PermissionHandler
	Role         *RoleHandler
	Menu         *MenuHandler
	Path         *PathHandler
	User         *UserHandler
	OperationLog *OperationLogHandler
	Migration    *MigrationHandler
}

type AuthHandler struct {
	svc *application.AuthService
}

func NewAuthHandler(svc *application.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// generateToken 生成 JWT Token
func generateToken(username string) (string, error) {
	jwtCfg := config.GetJWTConfig()
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Duration(jwtCfg.Expire) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtCfg.Secret))
}

// mockLogin 模拟登录，不依赖数据库
func mockLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if req.Username != "admin" || req.Password != "admin123" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "用户名或密码错误"})
		return
	}

	token, err := generateToken(req.Username)
	if err != nil {
		global.Logger.Error("生成 Token 失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "服务器错误"})
		return
	}

	permissions := []string{"user:list", "user:create", "user:edit", "user:delete"}

	global.Logger.Info("用户登录成功（模拟）", zap.String("username", req.Username))

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": LoginResponse{
			Token:       token,
			Permissions: permissions,
		},
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string   `json:"token"`
	Permissions []string `json:"permissions"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	if h.svc == nil {
		mockLogin(c)
		return
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	// 从数据库查询用户
	user, err := h.svc.GetUserByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "用户名或密码错误"})
			return
		}
		global.Logger.Error("查询用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "服务器错误"})
		return
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "用户名或密码错误"})
		return
	}

	// 生成 JWT Token
	token, err := generateToken(req.Username)
	if err != nil {
		global.Logger.Error("生成 Token 失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "服务器错误"})
		return
	}

	permissions := []string{"user:list", "user:create", "user:edit", "user:delete"}

	global.Logger.Info("用户登录成功", zap.String("username", req.Username))

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": LoginResponse{
			Token:       token,
			Permissions: permissions,
		},
	})
}

func (h *AuthHandler) GetPermissions(c *gin.Context) {
	// 从上下文中获取用户ID（由中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "未授权"})
		return
	}

	perms, err := h.svc.CheckPermission(userID.(uint), "")
	if err != nil {
		global.Logger.Warn("获取权限失败", zap.Any("userID", userID), zap.Error(err))
	}

	_ = perms
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": []model.Permission{},
	})
}
