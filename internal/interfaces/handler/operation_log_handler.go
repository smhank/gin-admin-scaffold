package handler

import (
	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/infras/global"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OperationLogHandler struct {
	db *gorm.DB
}

// mockOperationLogs 模拟操作日志数据
var mockOperationLogs = []model.OperationLog{
	{
		Model:     gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Username:  "admin",
		Action:    "用户登录",
		Method:    "POST",
		Path:      "/api/login",
		Params:    `{"username":"admin"}`,
		Result:    "success",
		Duration:  120,
		IP:        "127.0.0.1",
		UserAgent: "Mozilla/5.0",
	},
	{
		Model:     gorm.Model{ID: 2, CreatedAt: time.Now().Add(-time.Hour), UpdatedAt: time.Now().Add(-time.Hour)},
		Username:  "admin",
		Action:    "创建用户",
		Method:    "POST",
		Path:      "/api/users",
		Params:    `{"username":"test","nickname":"测试用户"}`,
		Result:    "success",
		Duration:  85,
		IP:        "127.0.0.1",
		UserAgent: "Mozilla/5.0",
	},
	{
		Model:     gorm.Model{ID: 3, CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now().Add(-2 * time.Hour)},
		Username:  "admin",
		Action:    "删除角色",
		Method:    "DELETE",
		Path:      "/api/roles/3",
		Params:    "",
		Result:    "fail",
		Duration:  30,
		IP:        "127.0.0.1",
		UserAgent: "Mozilla/5.0",
	},
}

var mockOpLogMu sync.Mutex
var mockOpLogNextID = uint(4)

func NewOperationLogHandler(db *gorm.DB) *OperationLogHandler {
	return &OperationLogHandler{db: db}
}

// ListOperationLogs 获取操作日志列表
func (h *OperationLogHandler) ListOperationLogs(c *gin.Context) {
	if h.db == nil {
		mockOpLogMu.Lock()
		defer mockOpLogMu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": mockOperationLogs})
		return
	}

	var logs []model.OperationLog
	query := h.db.Model(&model.OperationLog{})

	// 查询参数
	username := c.Query("username")
	action := c.Query("action")
	result := c.Query("result")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if result != "" {
		query = query.Where("result = ?", result)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	if err := query.Order("created_at desc").Find(&logs).Error; err != nil {
		global.Logger.Error("查询操作日志失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": logs})
}

// DeleteOperationLog 删除操作日志
func (h *OperationLogHandler) DeleteOperationLog(c *gin.Context) {
	id := c.Param("id")

	if h.db == nil {
		mockOpLogMu.Lock()
		defer mockOpLogMu.Unlock()
		for i, log := range mockOperationLogs {
			if id == "all" {
				mockOperationLogs = []model.OperationLog{}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "清空成功"})
				return
			}
			if log.ID == parseID(id) {
				mockOperationLogs = append(mockOperationLogs[:i], mockOperationLogs[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "日志不存在"})
		return
	}

	if id == "all" {
		if err := h.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.OperationLog{}).Error; err != nil {
			global.Logger.Error("清空操作日志失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "清空失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "清空成功"})
		return
	}

	if err := h.db.Delete(&model.OperationLog{}, id).Error; err != nil {
		global.Logger.Error("删除操作日志失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}

func parseID(id string) uint {
	var result uint
	for _, c := range id {
		if c >= '0' && c <= '9' {
			result = result*10 + uint(c-'0')
		} else {
			break
		}
	}
	return result
}
