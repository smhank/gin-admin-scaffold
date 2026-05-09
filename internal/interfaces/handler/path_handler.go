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

type PathHandler struct {
	db *gorm.DB
}

// mockPaths 模拟 API 路径数据
var mockPaths = []model.ApiPath{
	{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "登录", Path: "/api/login", Method: "POST", Desc: "用户登录"},
	{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "获取权限列表", Path: "/api/permissions", Method: "GET", Desc: "获取所有权限"},
	{Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "创建权限", Path: "/api/permissions", Method: "POST", Desc: "创建新权限"},
	{Model: gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "更新权限", Path: "/api/permissions/:id", Method: "PUT", Desc: "更新指定权限"},
	{Model: gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "删除权限", Path: "/api/permissions/:id", Method: "DELETE", Desc: "删除指定权限"},
	{Model: gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "获取角色列表", Path: "/api/roles", Method: "GET", Desc: "获取所有角色"},
	{Model: gorm.Model{ID: 7, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "创建角色", Path: "/api/roles", Method: "POST", Desc: "创建新角色"},
	{Model: gorm.Model{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "更新角色", Path: "/api/roles/:id", Method: "PUT", Desc: "更新指定角色"},
	{Model: gorm.Model{ID: 9, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "删除角色", Path: "/api/roles/:id", Method: "DELETE", Desc: "删除指定角色"},
	{Model: gorm.Model{ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "获取菜单树", Path: "/api/menus", Method: "GET", Desc: "获取菜单树结构"},
	{Model: gorm.Model{ID: 11, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "创建菜单", Path: "/api/menus", Method: "POST", Desc: "创建新菜单"},
	{Model: gorm.Model{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "更新菜单", Path: "/api/menus/:id", Method: "PUT", Desc: "更新指定菜单"},
	{Model: gorm.Model{ID: 13, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "删除菜单", Path: "/api/menus/:id", Method: "DELETE", Desc: "删除指定菜单"},
}

var mockPathMu sync.Mutex
var mockPathNextID = uint(14)

func NewPathHandler(db *gorm.DB) *PathHandler {
	return &PathHandler{db: db}
}

// ListPaths 获取 API 路径列表
func (h *PathHandler) ListPaths(c *gin.Context) {
	if h.db == nil {
		mockPathMu.Lock()
		defer mockPathMu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": mockPaths})
		return
	}
	var paths []model.ApiPath
	if err := h.db.Find(&paths).Error; err != nil {
		global.Logger.Error("查询API路径失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": paths})
}

type CreatePathRequest struct {
	Name   string `json:"name" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
	Desc   string `json:"desc"`
}

// CreatePath 创建 API 路径
func (h *PathHandler) CreatePath(c *gin.Context) {
	var req CreatePathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockPathMu.Lock()
		defer mockPathMu.Unlock()
		path := model.ApiPath{
			Model:  gorm.Model{ID: mockPathNextID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:   req.Name,
			Path:   req.Path,
			Method: req.Method,
			Desc:   req.Desc,
		}
		mockPathNextID++
		mockPaths = append(mockPaths, path)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": path})
		return
	}

	apiPath := model.ApiPath{Name: req.Name, Path: req.Path, Method: req.Method, Desc: req.Desc}
	if err := h.db.Create(&apiPath).Error; err != nil {
		global.Logger.Error("创建API路径失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": apiPath})
}

// UpdatePath 更新 API 路径
func (h *PathHandler) UpdatePath(c *gin.Context) {
	id := c.Param("id")
	var req CreatePathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockPathMu.Lock()
		defer mockPathMu.Unlock()
		for i, p := range mockPaths {
			if fmt.Sprintf("%d", p.ID) == id {
				mockPaths[i].Name = req.Name
				mockPaths[i].Path = req.Path
				mockPaths[i].Method = req.Method
				mockPaths[i].Desc = req.Desc
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功", "data": mockPaths[i]})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "API路径不存在"})
		return
	}

	if err := h.db.Model(&model.ApiPath{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name": req.Name, "path": req.Path, "method": req.Method, "desc": req.Desc,
	}).Error; err != nil {
		global.Logger.Error("更新API路径失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

// DeletePath 删除 API 路径
func (h *PathHandler) DeletePath(c *gin.Context) {
	id := c.Param("id")

	if h.db == nil {
		mockPathMu.Lock()
		defer mockPathMu.Unlock()
		for i, p := range mockPaths {
			if fmt.Sprintf("%d", p.ID) == id {
				mockPaths = append(mockPaths[:i], mockPaths[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "API路径不存在"})
		return
	}

	if err := h.db.Delete(&model.ApiPath{}, id).Error; err != nil {
		global.Logger.Error("删除API路径失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}
