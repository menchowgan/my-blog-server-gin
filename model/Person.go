package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Nickname string `gorm:"column:nickname;type:string;comment:'昵称'"`
	Password string `gorm:"column:password;type:string;comment:'密码'"`
	Hobbies  string `gorm:"column:hobbies;type:string;comment:'兴趣'"`
	Gender   string `gorm:"column:gender;type:varchar(1);comment:'性别'"`
	Fans     int16  `gorm:"column:fans;type:int(4);comment:'粉丝数量'"`
	Evaluate int16  `gorm:"column:evaluate;type:int(4);comment:'好评'"`
	Brief    string `gorm:"column:brief;type:string;comment:'简介'"`
	Avatar   string `gorm:"column:avatar;type:string;comment:'头像'"`
}

type UserEnroll struct {
	UserName string `json:"userName"`
	Passwrod string `json:"password"`
}
