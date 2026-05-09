package model

import "gorm.io/gorm"

// Migration 数据迁移记录
type Migration struct {
	gorm.Model
	Name      string `json:"name" gorm:"uniqueIndex:uni_migrations_name;not null;size:191;comment:迁移名称"`
	Batch     int    `json:"batch" gorm:"not null;default:1;comment:批次号"`
	AppliedAt string `json:"applied_at" gorm:"size:64;comment:应用时间"`
}

// TableName 指定表名
func (Migration) TableName() string {
	return "migrations"
}
