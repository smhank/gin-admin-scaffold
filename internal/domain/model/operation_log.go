package model

import "gorm.io/gorm"

type OperationLog struct {
	gorm.Model
	Username  string `json:"username" gorm:"size:64;not null;comment:操作人"`
	Action    string `json:"action" gorm:"size:128;not null;comment:操作动作"`
	Method    string `json:"method" gorm:"size:16;comment:请求方法"`
	Path      string `json:"path" gorm:"size:255;comment:请求路径"`
	Params    string `json:"params" gorm:"type:text;comment:请求参数"`
	Result    string `json:"result" gorm:"size:16;comment:操作结果 success/fail"`
	Duration  int64  `json:"duration" gorm:"comment:耗时(ms)"`
	IP        string `json:"ip" gorm:"size:64;comment:IP地址"`
	UserAgent string `json:"user_agent" gorm:"size:512;comment:User-Agent"`
}
