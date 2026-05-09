package model

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name      string `gorm:"not null;comment:菜单名称"`
	Icon      string `gorm:"default:'';comment:菜单图标"`
	Path      string `gorm:"default:'';comment:路由路径"`
	RouteName string `gorm:"default:'';comment:路由名称"`
	FilePath  string `gorm:"default:'';comment:文件路径"`
	ParentID  *uint  `gorm:"default:null;comment:父菜单ID"`
	Sort      int    `gorm:"default:0;comment:排序"`
	Status    int    `gorm:"default:1;comment:状态 1显示 0隐藏"`
	Roles     []Role `gorm:"many2many:role_menus;comment:可见角色"`
	Children  []Menu `gorm:"foreignKey:ParentID;comment:子菜单"`
}
