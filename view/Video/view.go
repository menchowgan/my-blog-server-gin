package video

import (
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"log"
)

func InsertVideoInfo(video *model.Video) error {
	dw := db.DB.GetDbW()

	log.Printf("UserId: %v, Avatar: %s, AudioUrl: %s", video.UserId, video.Avatar, video.VideoUrl)

	err := dw.Create(&video).Error
	if err == nil {
		return nil
	}
	return err
}
