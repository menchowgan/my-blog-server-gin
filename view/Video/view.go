package video

import (
	"gmc-blog-server/config"
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

func SearchByUserId(userid string) ([]VideoInfo, error) {
	dr := db.DB.GetDbR()
	var videos []model.Video
	err := dr.Where("userId = ?", userid).Order("created_at desc").Find(&videos).Error
	if err != nil {
		return []VideoInfo{}, err
	}
	videosArray := []VideoInfo{}
	for _, v := range videos {
		videosArray = append(videosArray, VideoInfo{
			ID:        v.ID,
			UserId:    int(v.UserId),
			Title:     v.Title,
			Artist:    v.Artist,
			Evalution: v.Evalution,
			VideoUrl:  config.VIDEO_QUERY_PATH + userid + "/" + v.VideoUrl,
			Avatar:    config.PHOTO_QUERY_PATH + userid + "/" + v.Avatar,
		})
	}
	return videosArray, nil
}
