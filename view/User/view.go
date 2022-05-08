package user

import (
	"gmc-blog-server/config"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"log"
)

func InsertUser(user model.User) (uint, error) {
	dw := db.DB.GetDbW()

	log.Printf("nickname: %s, gender: %s, avatar: %s", user.Nickname, user.Gender, user.Avatar)
	log.Println(user)

	err := dw.Create(&user).Error
	if err == nil {
		return user.ID, nil
	}
	return 0, err
}

func GerUserInfo(id string) PsersonSimpleIinfo {
	var user model.User
	dr := db.DB.GetDbR()

	dr.Where("id = ?", id).First(&user)

	u := PsersonSimpleIinfo{
		ID:       user.ID,
		Nickname: user.Nickname,
		Fans:     user.Fans,
		Avatar:   config.AVATAR_PATH + user.Avatar,
	}

	return u
}

func SearchUserBrief(id string) PersonInfoModel {
	var user model.User
	dr := db.DB.GetDbR()

	dr.Where("id = ?", id).First(&user)

	u := PersonInfoModel{
		ID:       user.ID,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Fans:     user.Fans,
		Brief:    user.Brief,
		Hobbies:  user.Hobbies,
		Avatar:   config.AVATAR_PATH + user.Avatar,
	}

	return u

}
