package photos

import "gorm.io/gorm"

type PhotoModel struct {
	gorm.Model
	UserId  uint   `gorm:"column:userId;type:bigint(20) unsigned;comment:'用户ID'"`
	ImgUrls string `gorm:"column:imgUrls;comment:'文章图片'"`
}
