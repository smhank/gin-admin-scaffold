package migration

import (
	"fmt"
	"gin-admin-base/internal/domain/model"

	"gorm.io/gorm"
)

// M20260507101614CreateTableTestMigration create_table_test_migration
type M20260507101614CreateTableTestMigration struct {
	BaseMigration
}

// NewM20260507101614CreateTableTestMigration 创建迁移
func NewM20260507101614CreateTableTestMigration() *M20260507101614CreateTableTestMigration {
	return &M20260507101614CreateTableTestMigration{
		BaseMigration: NewBaseMigration("m20260507_101614_create_table_test_migration"),
	}
}

// Up 执行迁移
func (m *M20260507101614CreateTableTestMigration) Up(db *gorm.DB) error {
	// TODO: 在这里编写迁移逻辑
	// 示例: db.Create(&model.Permission{Name: "示例权限", Code: "example:perm"})
	err := db.AutoMigrate(
		&model.TestMigrate{},
	)
	if err != nil {
		return fmt.Errorf("创建 test_migrates 表失败: %w", err)
	}
	fmt.Println("  ✓ test_migrates 表已创建")
	return nil
}

// Down 回滚迁移
func (m *M20260507101614CreateTableTestMigration) Down(db *gorm.DB) error {
	// TODO: 在这里编写回滚逻辑
	// 示例: db.Where("code = ?", "example:perm").Delete(&model.Permission{})
	err := db.AutoMigrate(
		&model.TestMigrate{},
	)
	if err != nil {
		return fmt.Errorf("删除 test_migrates 表失败: %w", err)
	}
	fmt.Println("  ✓ test_migrates 表已删除")
	return nil
}
