package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone    string `json:"phone" gorm:"uniqueIndex"`
	NickName string `json:"nick_name" gorm:"comment:姓名"`
	Avator   string `json:"avator" gorm:"comment:头像"`
	Sex      int    `json:"sex" gorm:"comment:性别"`
	Extra    string `json:"extra" gorm:"comment:额外信息"`
}
