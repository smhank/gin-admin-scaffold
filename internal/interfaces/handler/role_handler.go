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

// mockUsers 模拟用户数据
var mockUsers = []model.User{
	{
		Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Username: "admin",
		NickName: "Mr.奇森",
		Email:    "admin@example.com",
		Status:   1,
	},
	{
		Model:    gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Username: "a303176530",
		NickName: "用户1",
		Email:    "user1@example.com",
		Status:   1,
	},
}

var mockUserMu sync.Mutex
var mockUserNextID = uint(3)

type RoleHandler struct {
	db *gorm.DB
}

// mockRoles 模拟角色数据
var mockRoles = []model.Role{
	{
		Model:  gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:   "管理员",
		Status: 1,
		Permissions: []model.Permission{
			{Model: gorm.Model{ID: 1}, Name: "用户列表", Code: "user:list"},
			{Model: gorm.Model{ID: 2}, Name: "创建用户", Code: "user:create"},
			{Model: gorm.Model{ID: 3}, Name: "编辑用户", Code: "user:edit"},
			{Model: gorm.Model{ID: 4}, Name: "删除用户", Code: "user:delete"},
		},
	},
	{
		Model:  gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:   "普通用户",
		Status: 1,
		Permissions: []model.Permission{
			{Model: gorm.Model{ID: 1}, Name: "用户列表", Code: "user:list"},
		},
	},
}

var mockRoleMu sync.Mutex
var mockRoleNextID = uint(3)

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{db: db}
}

// ListRoles 获取角色列表
func (h *RoleHandler) ListRoles(c *gin.Context) {
	if h.db == nil {
		mockRoleMu.Lock()
		defer mockRoleMu.Unlock()
		// 为模拟数据添加Users字段
		type RoleWithUsers struct {
			model.Role
			Users []model.User `json:"Users"`
		}
		result := make([]RoleWithUsers, len(mockRoles))
		for i, r := range mockRoles {
			result[i] = RoleWithUsers{Role: r}
			// 从 mockRoles 中读取已分配的用户，而不是硬编码
			if len(mockRoles[i].Users) > 0 {
				result[i].Users = mockRoles[i].Users
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": result})
		return
	}
	var roles []model.Role
	if err := h.db.Preload("Permissions").Preload("Users").Find(&roles).Error; err != nil {
		global.Logger.Error("查询角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": roles})
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parentId"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	Permissions []uint `json:"permissions"`
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockRoleMu.Lock()
		defer mockRoleMu.Unlock()
		role := model.Role{
			Model:       gorm.Model{ID: mockRoleNextID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			ParentID:    req.ParentID,
			Sort:        req.Sort,
			Status:      req.Status,
		}
		for _, pid := range req.Permissions {
			for _, p := range mockPermissions {
				if p.ID == pid {
					role.Permissions = append(role.Permissions, p)
					break
				}
			}
		}
		mockRoleNextID++
		mockRoles = append(mockRoles, role)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": role})
		return
	}

	role := model.Role{Name: req.Name, Code: req.Code, Description: req.Description, ParentID: req.ParentID, Sort: req.Sort, Status: req.Status}
	if err := h.db.Create(&role).Error; err != nil {
		global.Logger.Error("创建角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "创建失败"})
		return
	}
	if len(req.Permissions) > 0 {
		var perms []model.Permission
		h.db.Find(&perms, req.Permissions)
		h.db.Model(&role).Association("Permissions").Replace(perms)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": role})
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockRoleMu.Lock()
		defer mockRoleMu.Unlock()
		for i, r := range mockRoles {
			if fmt.Sprintf("%d", r.ID) == id {
				mockRoles[i].Name = req.Name
				mockRoles[i].Code = req.Code
				mockRoles[i].Description = req.Description
				mockRoles[i].ParentID = req.ParentID
				mockRoles[i].Sort = req.Sort
				mockRoles[i].Status = req.Status
				mockRoles[i].Permissions = nil
				for _, pid := range req.Permissions {
					for _, p := range mockPermissions {
						if p.ID == pid {
							mockRoles[i].Permissions = append(mockRoles[i].Permissions, p)
							break
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功", "data": mockRoles[i]})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "角色不存在"})
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"code":        req.Code,
		"description": req.Description,
		"parent_id":   req.ParentID,
		"sort":        req.Sort,
		"status":      req.Status,
	}
	if err := h.db.Model(&model.Role{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		global.Logger.Error("更新角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "更新失败"})
		return
	}
	if len(req.Permissions) > 0 {
		var role model.Role
		h.db.First(&role, id)
		var perms []model.Permission
		h.db.Find(&perms, req.Permissions)
		h.db.Model(&role).Association("Permissions").Replace(perms)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	if h.db == nil {
		mockRoleMu.Lock()
		defer mockRoleMu.Unlock()
		for i, r := range mockRoles {
			if fmt.Sprintf("%d", r.ID) == id {
				mockRoles = append(mockRoles[:i], mockRoles[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "角色不存在"})
		return
	}

	if err := h.db.Delete(&model.Role{}, id).Error; err != nil {
		global.Logger.Error("删除角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}

// SetRolePermissions 设置角色权限
func (h *RoleHandler) SetRolePermissions(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Permissions []uint `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "设置成功"})
		return
	}

	var role model.Role
	if err := h.db.First(&role, id).Error; err != nil {
		global.Logger.Error("角色不存在", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "角色不存在"})
		return
	}

	if len(req.Permissions) > 0 {
		var perms []model.Permission
		h.db.Find(&perms, req.Permissions)
		h.db.Model(&role).Association("Permissions").Replace(perms)
	} else {
		h.db.Model(&role).Association("Permissions").Clear()
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "设置成功"})
}

// ListUsers 获取用户列表
func (h *RoleHandler) ListUsers(c *gin.Context) {
	if h.db == nil {
		mockUserMu.Lock()
		defer mockUserMu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": mockUsers})
		return
	}
	var users []model.User
	if err := h.db.Find(&users).Error; err != nil {
		global.Logger.Error("查询用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": users})
}

// AssignRoleToUser 分配角色给用户（全量覆盖角色的用户关联）
func (h *RoleHandler) AssignRoleToUser(c *gin.Context) {
	var req struct {
		RoleID  uint   `json:"roleId"`
		UserIDs []uint `json:"userIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockRoleMu.Lock()
		defer mockRoleMu.Unlock()
		// 更新 mock 角色的用户关联
		for i, r := range mockRoles {
			if r.ID == req.RoleID {
				var users []model.User
				for _, uid := range req.UserIDs {
					for _, u := range mockUsers {
						if u.ID == uid {
							users = append(users, u)
							break
						}
					}
				}
				mockRoles[i].Users = users
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "分配成功"})
		return
	}

	var role model.Role
	if err := h.db.First(&role, req.RoleID).Error; err != nil {
		global.Logger.Error("角色不存在", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "角色不存在"})
		return
	}

	// 获取所有用户
	var users []model.User
	if len(req.UserIDs) > 0 {
		if err := h.db.Find(&users, req.UserIDs).Error; err != nil {
			global.Logger.Error("用户不存在", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "用户不存在"})
			return
		}
	}

	// 从角色侧全量覆盖用户关联关系（不会影响用户的其他角色）
	if err := h.db.Model(&role).Association("Users").Replace(users); err != nil {
		global.Logger.Error("分配用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "分配失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "分配成功", "data": nil})
}
