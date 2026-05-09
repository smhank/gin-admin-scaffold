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

type MenuHandler struct {
	db *gorm.DB
}

// mockMenus 模拟菜单数据
var mockMenus = []model.Menu{
	{
		Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:     "首页",
		Icon:     "HomeFilled",
		Path:     "/dashboard",
		ParentID: nil,
		Sort:     1,
		Status:   1,
	},
	{
		Model:    gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:     "系统管理",
		Icon:     "Setting",
		Path:     "",
		ParentID: nil,
		Sort:     2,
		Status:   1,
		Children: []model.Menu{
			{
				Model:    gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Name:     "角色管理",
				Icon:     "UserFilled",
				Path:     "/roles",
				ParentID: uintPtr(2),
				Sort:     1,
				Status:   1,
			},
			{
				Model:    gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Name:     "权限管理",
				Icon:     "Lock",
				Path:     "/permissions",
				ParentID: uintPtr(2),
				Sort:     2,
				Status:   1,
			},
		},
	},
}

var mockMenuMu sync.Mutex
var mockMenuNextID = uint(5)

func uintPtr(v uint) *uint { return &v }

func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{db: db}
}

// ListMenus 获取菜单树
func (h *MenuHandler) ListMenus(c *gin.Context) {
	if h.db == nil {
		mockMenuMu.Lock()
		defer mockMenuMu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": mockMenus})
		return
	}
	var menus []model.Menu
	if err := h.db.Where("parent_id IS NULL").Preload("Children").Order("sort asc").Find(&menus).Error; err != nil {
		global.Logger.Error("查询菜单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": menus})
}

type CreateMenuRequest struct {
	Name      string `json:"name" binding:"required"`
	Icon      string `json:"icon"`
	Path      string `json:"path"`
	RouteName string `json:"routeName"`
	FilePath  string `json:"filePath"`
	ParentID  *uint  `json:"parentId"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status"`
}

// CreateMenu 创建菜单
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockMenuMu.Lock()
		defer mockMenuMu.Unlock()
		menu := model.Menu{
			Model:     gorm.Model{ID: mockMenuNextID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:      req.Name,
			Icon:      req.Icon,
			Path:      req.Path,
			RouteName: req.RouteName,
			FilePath:  req.FilePath,
			ParentID:  req.ParentID,
			Sort:      req.Sort,
			Status:    req.Status,
		}
		mockMenuNextID++
		if req.ParentID == nil {
			mockMenus = append(mockMenus, menu)
		} else {
			addChildMenu(mockMenus, *req.ParentID, menu)
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": menu})
		return
	}

	menu := model.Menu{
		Name: req.Name, Icon: req.Icon, Path: req.Path,
		RouteName: req.RouteName, FilePath: req.FilePath,
		ParentID: req.ParentID, Sort: req.Sort, Status: req.Status,
	}
	if err := h.db.Create(&menu).Error; err != nil {
		global.Logger.Error("创建菜单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": menu})
}

func addChildMenu(menus []model.Menu, parentID uint, child model.Menu) {
	for i, m := range menus {
		if m.ID == parentID {
			menus[i].Children = append(menus[i].Children, child)
			return
		}
		if len(m.Children) > 0 {
			addChildMenu(m.Children, parentID, child)
		}
	}
}

// UpdateMenu 更新菜单
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		mockMenuMu.Lock()
		defer mockMenuMu.Unlock()
		if updated := updateMockMenu(mockMenus, id, req); updated {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "菜单不存在"})
		return
	}

	updates := map[string]interface{}{
		"name": req.Name, "icon": req.Icon, "path": req.Path,
		"route_name": req.RouteName, "file_path": req.FilePath,
		"parent_id": req.ParentID, "sort": req.Sort, "status": req.Status,
	}
	if err := h.db.Model(&model.Menu{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		global.Logger.Error("更新菜单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

func updateMockMenu(menus []model.Menu, id string, req CreateMenuRequest) bool {
	for i, m := range menus {
		if fmt.Sprintf("%d", m.ID) == id {
			menus[i].Name = req.Name
			menus[i].Icon = req.Icon
			menus[i].Path = req.Path
			menus[i].RouteName = req.RouteName
			menus[i].FilePath = req.FilePath
			menus[i].ParentID = req.ParentID
			menus[i].Sort = req.Sort
			menus[i].Status = req.Status
			return true
		}
		if len(m.Children) > 0 {
			if updateMockMenu(m.Children, id, req) {
				return true
			}
		}
	}
	return false
}

// DeleteMenu 删除菜单
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	if h.db == nil {
		mockMenuMu.Lock()
		defer mockMenuMu.Unlock()
		if deleted := deleteMockMenu(&mockMenus, id); deleted {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "菜单不存在"})
		return
	}

	if err := h.db.Delete(&model.Menu{}, id).Error; err != nil {
		global.Logger.Error("删除菜单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}

// SetMenuRoles 设置菜单可见角色
func (h *MenuHandler) SetMenuRoles(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		RoleIDs []uint `json:"roleIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	if h.db == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "设置成功"})
		return
	}

	var menu model.Menu
	if err := h.db.First(&menu, id).Error; err != nil {
		global.Logger.Error("菜单不存在", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "菜单不存在"})
		return
	}

	if len(req.RoleIDs) > 0 {
		var roles []model.Role
		h.db.Find(&roles, req.RoleIDs)
		h.db.Model(&menu).Association("Roles").Replace(roles)
	} else {
		h.db.Model(&menu).Association("Roles").Clear()
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "设置成功"})
}

func deleteMockMenu(menus *[]model.Menu, id string) bool {
	for i, m := range *menus {
		if fmt.Sprintf("%d", m.ID) == id {
			*menus = append((*menus)[:i], (*menus)[i+1:]...)
			return true
		}
		if len(m.Children) > 0 {
			if deleteMockMenu(&m.Children, id) {
				(*menus)[i].Children = m.Children
				return true
			}
		}
	}
	return false
}
