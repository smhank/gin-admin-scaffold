package handler

import (
	"fmt"
	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/infras/global"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PermissionHandler struct {
	db *gorm.DB
}

// mockPermissions 模拟权限数据
var mockPermissions = []model.Permission{
	{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "用户列表", Code: "user:list"},
	{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "创建用户", Code: "user:create"},
	{Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "编辑用户", Code: "user:edit"},
	{Model: gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "删除用户", Code: "user:delete"},
	{Model: gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "角色列表", Code: "role:list"},
	{Model: gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "创建角色", Code: "role:create"},
	{Model: gorm.Model{ID: 7, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "编辑角色", Code: "role:edit"},
	{Model: gorm.Model{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "删除角色", Code: "role:delete"},
	{Model: gorm.Model{ID: 9, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "权限列表", Code: "permission:list"},
	{Model: gorm.Model{ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "创建权限", Code: "permission:create"},
	{Model: gorm.Model{ID: 11, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "编辑权限", Code: "permission:edit"},
	{Model: gorm.Model{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "删除权限", Code: "permission:delete"},
}

var mockPermMu sync.Mutex
var mockPermNextID = uint(13)

func NewPermissionHandler(db *gorm.DB) *PermissionHandler {
	return &PermissionHandler{db: db}
}

// ListPermissions 获取权限列表
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	if h.db == nil {
		// 模拟数据
		mockPermMu.Lock()
		defer mockPermMu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": mockPermissions})
		return
	}
	var perms []model.Permission
	if err := h.db.Find(&perms).Error; err != nil {
		global.Logger.Error("查询权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": perms})
}

type CreatePermissionRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockPermMu.Lock()
		defer mockPermMu.Unlock()
		perm := model.Permission{
			Model: gorm.Model{ID: mockPermNextID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:  req.Name,
			Code:  req.Code,
		}
		mockPermNextID++
		mockPermissions = append(mockPermissions, perm)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": perm})
		return
	}

	perm := model.Permission{Name: req.Name, Code: req.Code}
	if err := h.db.Create(&perm).Error; err != nil {
		global.Logger.Error("创建权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": perm})
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockPermMu.Lock()
		defer mockPermMu.Unlock()
		for i, p := range mockPermissions {
			if fmt.Sprintf("%d", p.ID) == id {
				mockPermissions[i].Name = req.Name
				mockPermissions[i].Code = req.Code
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功", "data": mockPermissions[i]})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "权限不存在"})
		return
	}

	if err := h.db.Model(&model.Permission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name": req.Name, "code": req.Code,
	}).Error; err != nil {
		global.Logger.Error("更新权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")

	if h.db == nil {
		mockPermMu.Lock()
		defer mockPermMu.Unlock()
		for i, p := range mockPermissions {
			if fmt.Sprintf("%d", p.ID) == id {
				mockPermissions = append(mockPermissions[:i], mockPermissions[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "权限不存在"})
		return
	}

	if err := h.db.Delete(&model.Permission{}, id).Error; err != nil {
		global.Logger.Error("删除权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}
