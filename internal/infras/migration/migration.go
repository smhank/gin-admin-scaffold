package migration

import "gorm.io/gorm"

// Migration 迁移接口（类似 Yii2 的 Migration 基类）
type Migration interface {
	// Name 返回迁移名称（如 "m20250506_000001_init_seed_data"）
	Name() string
	// Up 执行迁移
	Up(db *gorm.DB) error
	// Down 回滚迁移
	Down(db *gorm.DB) error
}

// BaseMigration 基础迁移（类似 Yii2 的 yii\db\Migration）
type BaseMigration struct {
	name string
}

// NewBaseMigration 创建基础迁移
func NewBaseMigration(name string) BaseMigration {
	return BaseMigration{name: name}
}

// Name 返回迁移名称
func (m *BaseMigration) Name() string {
	return m.name
}

// Down 默认回滚方法（子类可覆盖）
func (m *BaseMigration) Down(db *gorm.DB) error {
	return nil
}
