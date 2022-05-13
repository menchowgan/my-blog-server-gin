package model

import "gorm.io/gorm"

type Music struct {
	gorm.Model
	UserId    uint   `gorm:"column:userId;type:bigint(20) unsigned;comment:'用户ID'"`
	Avatar    string `gorm:"column:avatar;type:string;comment:'封面'"`
	Title     string `gorm:"column:title;type:string;comment:'歌名'"`
	Artist    string `gorm:"column:artist;type:string;comment:'创作者'"`
	Evalution string `gorm:"column:evalution;type:string;comment:'评价'"`
	AudioUrl  string `gorm:"column:audioUrl;comment:''"`
}
