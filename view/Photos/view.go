package photos

import (
	"errors"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"log"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func PhotoUpdate(fileName string, userId string) error {
	photosStringModel, err := PhotosQueryByUserId(userId)
	if err != nil {
		return err
	}
	log.Println(photosStringModel.UserId, photosStringModel.ID, photosStringModel.ImgUrls)

	if photosStringModel.ID == 0 {
		i, _ := strconv.Atoi(userId)

		photoInfo := model.Photos{
			ImgUrls: fileName,
			UserId:  uint(i),
		}
		log.Printf("photo info is: %s\\ %v", photoInfo.ImgUrls, photoInfo.UserId)
		dw := db.DB.GetDbW()
		err := dw.Create(&photoInfo).Error
		if err != nil {
			return err
		}
	} else {
		log.Println("new file image: ", fileName)
		photosStringModel.ImgUrls = photosStringModel.ImgUrls + ";" + fileName
		dw := db.DB.GetDbW()
		err = dw.Save(&photosStringModel).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func PhotosQueryByUserId(userId string) (model.Photos, error) {
	dr := db.DB.GetDbR()

	log.Println("user id: ", userId)

	var photos model.Photos
	err := dr.Where("userId = ?", userId).First(&photos).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("未找到数据")
		} else {
			return model.Photos{}, err
		}
	}
	log.Println("photos urls: ", userId, photos.ImgUrls)

	return photos, nil
}

func PhotoDeleteByFileName(uid string, filename string) error {
	dr := db.DB.GetDbR()
	dw := db.DB.GetDbW()

	var photo model.Photos
	err := dr.Where("userId = ?", uid).First(&photo).Error

	log.Println("photos query: ", photo)

	if err != nil {
		return err
	}

	if photo.ID != 0 && photo.ImgUrls != "" {
		imgUrls := strings.Split(photo.ImgUrls, ";")
		imgs := []string{}
		for _, url := range imgUrls {
			if url != filename {
				imgs = append(imgs, url)
			}
		}
		photo.ImgUrls = strings.Join(imgs, ";")
		log.Println("new photos uri: ", photo.ImgUrls)
		err = dw.Save(&photo).Error
		if err != nil {
			return err
		}
	}

	return nil
}
