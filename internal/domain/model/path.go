package model

import "gorm.io/gorm"

type ApiPath struct {
	gorm.Model
	Name   string `gorm:"not null;comment:接口名称"`
	Path   string `gorm:"not null;comment:接口路径"`
	Method string `gorm:"default:GET;comment:请求方法"`
	Desc   string `gorm:"default:'';comment:接口描述"`
}
