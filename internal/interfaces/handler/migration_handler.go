package handler

import (
	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/infras/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MigrationHandler struct {
	db *gorm.DB
}

func NewMigrationHandler(db *gorm.DB) *MigrationHandler {
	return &MigrationHandler{db: db}
}

// ListMigrations 获取迁移记录列表
func (h *MigrationHandler) ListMigrations(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": []model.Migration{}})
		return
	}

	var migrations []model.Migration
	if err := h.db.Model(&model.Migration{}).Order("batch desc, created_at desc").Find(&migrations).Error; err != nil {
		global.Logger.Error("查询迁移记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": migrations})
}
