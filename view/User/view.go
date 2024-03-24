package user

import (
	"encoding/json"
	"gmc-blog-server/config"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	myRedis "gmc-blog-server/redis"
	article "gmc-blog-server/view/Article"
	music "gmc-blog-server/view/Music"
	photos "gmc-blog-server/view/Photos"
	video "gmc-blog-server/view/Video"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func InsertUser(user *model.User) error {
	dw := db.DB.GetDbW()

	log.Printf("nickname: %s, gender: %s, avatar: %s", user.Nickname, user.Gender, user.Avatar)
	log.Println(user)

	err := dw.Create(&user).Error
	if err == nil {
		return nil
	}
	return err
}

func Save(user *model.User) (*model.User, error) {
	dw := db.DB.GetDbW()
	err := dw.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"nickname": user.Nickname,
		"hobbies":  user.Hobbies,
		"gender":   user.Gender,
		"brief":    user.Brief,
		"avatar":   user.Avatar,
	}).Error
	if err != nil {
		return nil, err
	}
	searchAndUpdateCache(strconv.Itoa(int(user.ID)))
	return user, nil
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
		for idx, url := range urls {
			model := photos.PhotoInfo{
				ID:  uint(idx + 1),
				Url: config.PHOTO_QUERY_PATH + id + "/" + url,
			}
			imgUrls = append(imgUrls, model)
		}
	}

	articles, err := article.ArticleSimplaeInfosQueryByUserId(id)
	if err != nil {
		return PsersonSimpleIinfo{}, err
	}

	articleSIs := getArticleSimpleInfo(id, articles)

	u := PsersonSimpleIinfo{
		ID:                 user.ID,
		Nickname:           user.Nickname,
		Fans:               user.Fans,
		Photos:             imgUrls,
		ArticleSimpleInfos: articleSIs,
		Avatar:             config.PHOTO_QUERY_PATH + strconv.Itoa(int(user.ID)) + "/" + user.Avatar,
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

	audioArray, err := music.SearchByUserId(id)
	if err != nil {
		return PersonInfoModel{}, err
	}

	videosArray, err := video.SearchByUserId(id)
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
		Audios:   audioArray,
		Videos:   videosArray,
		Avatar:   config.PHOTO_QUERY_PATH + id + "/" + user.Avatar,
	}

	return u, nil
}

func getArticleSimpleInfo(id string, articles []model.Articles) []article.ArticleSimpleInfoModel {
	var articleSIs []article.ArticleSimpleInfoModel //id, userId, imgUrl, title, content, created_at
	if len(articles) > 0 {
		for _, a := range articles {
			articleSIs = append(articleSIs, article.ArticleSimpleInfoModel{
				ID:     int64(a.ID),
				Title:  a.Title,
				UserId: a.UserId,
				ImgUrl: config.PHOTO_QUERY_PATH + id + "/article/" + a.ImgUrl,
				Brief:  a.Brief,
				Date:   a.CreatedAt,
			})
		}
	}

	return articleSIs
}

func GetUserAllInfo(id string) (UserEnrollModel, error) {
	var user model.User
	dr := db.DB.GetDbR()

	err := dr.Where("id = ?", id).First(&user).Error
	if err != nil {
		return UserEnrollModel{}, err
	}
	return UserEnrollModel{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Nickname:  user.Nickname,
		Gender:    user.Gender,
		Hobbies:   user.Hobbies,
		Evaluate:  int(user.Evaluate),
		Brief:     user.Brief,
		Avatar:    config.PHOTO_QUERY_PATH + id + "/" + user.Avatar,
	}, nil
}

func FindUser(id string) (model.User, error) {
	s, err := myRedis.RedisDb.Get(myRedis.GetKey("user", id)).Result()
	log.Println("cache searched", s, err)
	var user model.User

	if err != nil || err == redis.Nil {
		return searchAndUpdateCache(id)
	} else {
		if s != "" {
			log.Println("cached user", s)
			if err := json.Unmarshal([]byte(s), &user); err == nil {
				return user, nil
			}
		}
		return user, nil
	}

}

func searchAndUpdateCache(id string) (model.User, error) {
	dr := db.DB.GetDbR()
	var user model.User

	err := dr.Select("id, nickname, fans, avatar").Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	log.Println("user info: ", user)
	if b, err := json.Marshal(&user); err == nil {
		_, err = myRedis.RedisDb.Set(myRedis.GetKey("user", id), b, 96*time.Hour).Result()
	}

	return user, nil
}
