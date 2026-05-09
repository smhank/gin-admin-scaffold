package migration

import (
	"fmt"

	"gin-admin-base/internal/domain/model"

	"gorm.io/gorm"
)

// M20250506000002AddOperationLogMenu 添加操作历史菜单
type M20250506000002AddOperationLogMenu struct {
	BaseMigration
}

// NewM20250506000002AddOperationLogMenu 创建迁移
func NewM20250506000002AddOperationLogMenu() *M20250506000002AddOperationLogMenu {
	return &M20250506000002AddOperationLogMenu{
		BaseMigration: NewBaseMigration("m20250506_000002_add_operation_log_menu"),
	}
}

// Up 执行迁移
func (m *M20250506000002AddOperationLogMenu) Up(db *gorm.DB) error {
	// 检查菜单是否已存在
	var count int64
	db.Model(&model.Menu{}).Where("path = ?", "/operation-logs").Count(&count)
	if count > 0 {
		fmt.Println("  操作历史菜单已存在，跳过")
		return nil
	}

	// 获取系统管理菜单
	var systemMenu model.Menu
	if err := db.Where("path = ?", "/system").First(&systemMenu).Error; err != nil {
		return fmt.Errorf("查询系统管理菜单失败: %w", err)
	}

	// 创建操作历史菜单
	menu := model.Menu{
		Name:     "操作历史",
		Icon:     "Document",
		Path:     "/operation-logs",
		ParentID: &systemMenu.ID,
		Sort:     5,
		Status:   1,
	}
	if err := db.Create(&menu).Error; err != nil {
		return fmt.Errorf("创建操作历史菜单失败: %w", err)
	}
	fmt.Println("  ✓ 操作历史菜单已创建")

	return nil
}

// Down 回滚迁移
func (m *M20250506000002AddOperationLogMenu) Down(db *gorm.DB) error {
	if err := db.Where("path = ?", "/operation-logs").Delete(&model.Menu{}).Error; err != nil {
		return fmt.Errorf("删除操作历史菜单失败: %w", err)
	}
	fmt.Println("  ✓ 操作历史菜单已删除")
	return nil
}
