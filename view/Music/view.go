package music

import (
	"errors"
	"gmc-blog-server/config"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"log"

	"gorm.io/gorm"
)

func InsertMusicInfo(audio model.Music) (uint, error) {
	dw := db.DB.GetDbW()

	log.Printf("UserId: %v, Avatar: %s, AudioUrl: %s", audio.UserId, audio.Avatar, audio.AudioUrl)

	err := dw.Create(&audio).Error
	if err == nil {
		return audio.ID, nil
	}
	return 0, err
}

func SearchByUserId(userid string) ([]MusicInfo, error) {
	dr := db.DB.GetDbR()
	var audios []model.Music
	err := dr.Where("userId = ?", userid).Order("created_at desc").Find(&audios).Error
	if err != nil {
		return []MusicInfo{}, err
	}
	audioArray := []MusicInfo{}
	for _, audio := range audios {
		audioArray = append(audioArray, MusicInfo{
			ID:        audio.ID,
			UserId:    int(audio.UserId),
			Title:     audio.Title,
			Artist:    audio.Artist,
			Evalution: audio.Evalution,
			Avatar:    config.PHOTO_QUERY_PATH + userid + "/" + audio.Avatar,
			AudioUrl:  config.MUSCI_QUERY_PATH + userid + "/" + audio.AudioUrl,
		})
	}
	return audioArray, nil
}

func MusicQueryByUserIdSimplaeLife(userId string) ([]model.Music, error) {
	dr := db.DB.GetDbR()

	log.Println("user id: ", userId)

	var musicSI []model.Music
	err := dr.Where("userId = ?", userId).Order("created_at desc").Find(&musicSI).Limit(5).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("未找到数据")
		} else {
			return []model.Music{}, err
		}
	}

	return musicSI, nil
}
