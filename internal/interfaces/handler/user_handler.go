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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

// mockUserData 模拟用户数据
var mockUserData = []model.User{
	{
		Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Username: "admin",
		NickName: "Mr.奇森",
		Email:    "admin@example.com",
		Phone:    "13800138000",
		Status:   1,
	},
	{
		Model:    gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Username: "a303176530",
		NickName: "用户1",
		Email:    "user1@example.com",
		Phone:    "13800138001",
		Status:   1,
	},
}

var mockUserDataMu sync.Mutex
var mockUserDataNextID = uint(3)

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	if h.db == nil {
		mockUserDataMu.Lock()
		defer mockUserDataMu.Unlock()
		// 不返回密码
		type SafeUser struct {
			ID        uint             `json:"ID"`
			CreatedAt time.Time        `json:"CreatedAt"`
			UpdatedAt time.Time        `json:"UpdatedAt"`
			Username  string           `json:"Username"`
			NickName  string           `json:"NickName"`
			Phone     string           `json:"Phone"`
			Email     string           `json:"Email"`
			Avatar    string           `json:"Avatar"`
			Status    model.UserStatus `json:"Status"`
			Roles     []model.Role     `json:"Roles"`
		}
		result := make([]SafeUser, len(mockUserData))
		for i, u := range mockUserData {
			result[i] = SafeUser{
				ID:        u.ID,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
				Username:  u.Username,
				NickName:  u.NickName,
				Phone:     u.Phone,
				Email:     u.Email,
				Avatar:    u.Avatar,
				Status:    u.Status,
				Roles:     u.Roles,
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": result})
		return
	}

	var users []model.User
	if err := h.db.Preload("Roles").Find(&users).Error; err != nil {
		global.Logger.Error("查询用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": users})
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockUserDataMu.Lock()
		defer mockUserDataMu.Unlock()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user := model.User{
			Model:    gorm.Model{ID: mockUserDataNextID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Username: req.Username,
			Password: string(hashedPassword),
			NickName: req.NickName,
			Phone:    req.Phone,
			Email:    req.Email,
			Status:   model.UserStatus(req.Status),
		}
		if user.Status == 0 {
			user.Status = model.UserStatusActive
		}
		// 分配角色
		for _, rid := range req.RoleIDs {
			for _, r := range mockRoles {
				if r.ID == rid {
					user.Roles = append(user.Roles, r)
					break
				}
			}
		}
		mockUserDataNextID++
		mockUserData = append(mockUserData, user)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": user})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Error("密码加密失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "创建失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		NickName: req.NickName,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   model.UserStatus(req.Status),
	}
	if user.Status == 0 {
		user.Status = model.UserStatusActive
	}

	if err := h.db.Create(&user).Error; err != nil {
		global.Logger.Error("创建用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "创建失败"})
		return
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		var roles []model.Role
		h.db.Find(&roles, req.RoleIDs)
		h.db.Model(&user).Association("Roles").Replace(roles)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": user})
}

type UpdateUserRequest struct {
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockUserDataMu.Lock()
		defer mockUserDataMu.Unlock()
		for i, u := range mockUserData {
			if fmt.Sprintf("%d", u.ID) == id {
				mockUserData[i].NickName = req.NickName
				mockUserData[i].Phone = req.Phone
				mockUserData[i].Email = req.Email
				if req.Status != 0 {
					mockUserData[i].Status = model.UserStatus(req.Status)
				}
				// 更新角色
				if req.RoleIDs != nil {
					mockUserData[i].Roles = nil
					for _, rid := range req.RoleIDs {
						for _, r := range mockRoles {
							if r.ID == rid {
								mockUserData[i].Roles = append(mockUserData[i].Roles, r)
								break
							}
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功", "data": mockUserData[i]})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	updates := map[string]interface{}{
		"nick_name": req.NickName,
		"phone":     req.Phone,
		"email":     req.Email,
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := h.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		global.Logger.Error("更新用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "更新失败"})
		return
	}

	// 更新角色
	if req.RoleIDs != nil {
		var user model.User
		h.db.First(&user, id)
		var roles []model.Role
		if len(req.RoleIDs) > 0 {
			h.db.Find(&roles, req.RoleIDs)
		}
		h.db.Model(&user).Association("Roles").Replace(roles)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if h.db == nil {
		mockUserDataMu.Lock()
		defer mockUserDataMu.Unlock()
		for i, u := range mockUserData {
			if fmt.Sprintf("%d", u.ID) == id {
				mockUserData = append(mockUserData[:i], mockUserData[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	if err := h.db.Delete(&model.User{}, id).Error; err != nil {
		global.Logger.Error("删除用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockUserDataMu.Lock()
		defer mockUserDataMu.Unlock()
		for i, u := range mockUserData {
			if fmt.Sprintf("%d", u.ID) == id {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
				mockUserData[i].Password = string(hashedPassword)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "密码重置成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Error("密码加密失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "重置失败"})
		return
	}

	if err := h.db.Model(&model.User{}).Where("id = ?", id).Update("password", string(hashedPassword)).Error; err != nil {
		global.Logger.Error("重置密码失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "重置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "密码重置成功"})
}

// UpdateStatus 更新用户状态
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockUserDataMu.Lock()
		defer mockUserDataMu.Unlock()
		for i, u := range mockUserData {
			if fmt.Sprintf("%d", u.ID) == id {
				mockUserData[i].Status = model.UserStatus(req.Status)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "状态更新成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	if err := h.db.Model(&model.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		global.Logger.Error("更新状态失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "状态更新成功"})
}

// AssignRoles 分配角色给用户
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		RoleIDs []uint `json:"roleIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockUserDataMu.Lock()
		defer mockUserDataMu.Unlock()
		for i, u := range mockUserData {
			if fmt.Sprintf("%d", u.ID) == id {
				mockUserData[i].Roles = nil
				for _, rid := range req.RoleIDs {
					for _, r := range mockRoles {
						if r.ID == rid {
							mockUserData[i].Roles = append(mockUserData[i].Roles, r)
							break
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "分配成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	var user model.User
	if err := h.db.First(&user, id).Error; err != nil {
		global.Logger.Error("用户不存在", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	var roles []model.Role
	if len(req.RoleIDs) > 0 {
		if err := h.db.Find(&roles, req.RoleIDs).Error; err != nil {
			global.Logger.Error("角色不存在", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "角色不存在"})
			return
		}
	}

	if err := h.db.Model(&user).Association("Roles").Replace(roles); err != nil {
		global.Logger.Error("分配角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "分配失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "分配成功"})
}
