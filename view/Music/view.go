package music

import (
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"log"
)

func MusicUpdate(fileName string, userId string) error {
	return nil
}

func InsertMusicInfo(audio model.Music) (uint, error) {
	dw := db.DB.GetDbW()

	log.Printf("UserId: %v, Avatar: %s, AudioUrl: %s", audio.UserId, audio.Avatar, audio.AudioUrl)

	err := dw.Create(&audio).Error
	if err == nil {
		return audio.ID, nil
	}
	return 0, err
}
