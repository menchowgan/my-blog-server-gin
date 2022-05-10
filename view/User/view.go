package user

import (
	"gmc-blog-server/config"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	photos "gmc-blog-server/view/Photos"
	"log"
	"strconv"
	"strings"
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

func GerUserInfo(id string) (PsersonSimpleIinfo, error) {
	var user model.User
	var photo model.Photos
	dr := db.DB.GetDbR()

	err := dr.Select("id, nickname, fans, avatar").Where("id = ?", id).First(&user).Error
	if err != nil {
		return PsersonSimpleIinfo{}, err
	}
	log.Println("user info: ", user)

	photo, err = photos.PhotosQueryByUserId(id)
	if err != nil {
		return PsersonSimpleIinfo{}, err
	}

	imgUrls := []photos.PhotoInfo{}

	if photo.ImgUrls != "" && len(photo.ImgUrls) > 0 {
		urls := strings.Split(photo.ImgUrls, ";")
		for _, url := range urls {
			model := photos.PhotoInfo{
				ID:  photo.ID,
				Url: config.PHOTO_QUERY_PATH + id + "/" + url,
			}
			imgUrls = append(imgUrls, model)
		}
	}

	u := PsersonSimpleIinfo{
		ID:       user.ID,
		Nickname: user.Nickname,
		Fans:     user.Fans,
		Photos:   imgUrls,
		Avatar:   config.PHOTO_QUERY_PATH + strconv.Itoa(int(user.ID)) + "/" + user.Avatar,
	}

	return u, nil
}

func SearchUserBrief(id string) (PersonInfoModel, error) {
	var user model.User
	dr := db.DB.GetDbR()

	err := dr.Where("id = ?", id).First(&user).Error
	if err != nil {
		return PersonInfoModel{}, err
	}

	u := PersonInfoModel{
		ID:       user.ID,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Fans:     user.Fans,
		Brief:    user.Brief,
		Hobbies:  user.Hobbies,
		Avatar:   config.PHOTO_QUERY_PATH + strconv.Itoa(int(user.ID)) + "/" + user.Avatar,
	}

	return u, nil

}
