package model

import (
	"gorm.io/gorm"
)

type Articles struct {
	gorm.Model
	UserId  uint   `gorm:"column:userId;type:bigint(20) unsigned;comment:'用户ID'"`
	ImgUrl  string `gorm:"column:imgUrl;comment:'文章图片'"`
	Title   string `gorm:"column:title;not null;type:varchar(100);comment:'文章标题'"`
	Content string `gorm:"column:content;type:string;comment:'内容'"`
	Brief   string `gorm:"column:brief;type:string;comment:'摘要'"`
	Type    string `gorm:"column:type;type:varchar(100);comment:'类型'"`
}
