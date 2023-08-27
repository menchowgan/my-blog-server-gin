package model

import "gorm.io/gorm"

type Plan struct {
	gorm.Model

	UserId  uint   `gorm:"column:userId;type:bigint(20) unsigned;comment:'用户ID'"`
	Content string `gorm:"volumn:userId;type:string;comment:'内容'"`
}
