package migration

import (
	"fmt"

	"gin-admin-base/internal/domain/model"

	"gorm.io/gorm"
)

// M20250506000000InitTables 初始化数据库表结构
type M20250506000000InitTables struct {
	BaseMigration
}

// NewM20250506000000InitTables 创建迁移
func NewM20250506000000InitTables() *M20250506000000InitTables {
	return &M20250506000000InitTables{
		BaseMigration: NewBaseMigration("m20250506_000000_init_tables"),
	}
}

// Up 执行迁移
func (m *M20250506000000InitTables) Up(db *gorm.DB) error {
	// 先创建 migrations 表（迁移系统自身需要）
	if err := db.AutoMigrate(&model.Migration{}); err != nil {
		return fmt.Errorf("创建 migrations 表失败: %w", err)
	}
	fmt.Println("  ✓ migrations 表已创建")

	// 创建所有业务表
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.Menu{},
		&model.ApiPath{},
		&model.OperationLog{},
	)
	if err != nil {
		return fmt.Errorf("创建业务表失败: %w", err)
	}
	fmt.Println("  ✓ 业务表已创建")

	return nil
}

// Down 回滚迁移
func (m *M20250506000000InitTables) Down(db *gorm.DB) error {
	// 删除所有业务表
	tables := []interface{}{
		&model.OperationLog{},
		&model.ApiPath{},
		&model.Menu{},
		&model.Permission{},
		&model.Role{},
		&model.User{},
	}
	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("删除表失败: %w", err)
		}
	}
	fmt.Println("  ✓ 业务表已删除")
	return nil
}
