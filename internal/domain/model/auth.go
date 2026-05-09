package model

import "gorm.io/gorm"

type UserStatus int

const (
	UserStatusActive   UserStatus = 1
	UserStatusDisabled UserStatus = 0
)

type User struct {
	gorm.Model
	Username string     `gorm:"unique;not null;comment:用户名"`
	Password string     `gorm:"not null;comment:密码"`
	NickName string     `gorm:"default:'';comment:昵称"`
	Phone    string     `gorm:"default:'';comment:手机号"`
	Email    string     `gorm:"default:'';comment:邮箱"`
	Avatar   string     `gorm:"default:'';comment:头像"`
	Status   UserStatus `gorm:"default:1;comment:状态 1正常 0禁用"`
	Roles    []Role     `gorm:"many2many:user_roles;comment:角色"`
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique;not null;comment:角色名称"`
	Code        string       `gorm:"default:'';comment:角色编码"`
	Description string       `gorm:"default:'';comment:角色描述"`
	ParentID    *uint        `gorm:"default:null;comment:上级角色ID"`
	Sort        int          `gorm:"default:0;comment:排序"`
	Status      int          `gorm:"default:1;comment:状态 1正常 0禁用"`
	Permissions []Permission `gorm:"many2many:role_permissions;comment:权限"`
	Users       []User       `gorm:"many2many:user_roles;comment:用户"`
	Children    []Role       `gorm:"foreignkey:ParentID;comment:子角色"`
}

type Permission struct {
	gorm.Model
	Name string `gorm:"unique;not null;comment:权限名称"`
	Code string `gorm:"unique;not null;comment:权限标识"`
}
