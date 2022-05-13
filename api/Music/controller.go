package music

import (
	fileapi "gmc-blog-server/api/File"
	"gmc-blog-server/config"
	"gmc-blog-server/model"
	music "gmc-blog-server/view/Music"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserMusicUpload(c *gin.Context) error {
	var audio model.Music

	if err := c.ShouldBind(&audio); err != nil {
		return err
	}
	log.Println("audio post: ", audio.AudioUrl)
	log.Println("audio post: ", audio.UserId)
	log.Println("audio post: ", audio.Avatar)
	log.Println("audio post: ", audio.Title)
	log.Println("audio post: ", audio.Artist)
	log.Println("audio post: ", audio.Evalution)

	audioId, err := music.InsertMusicInfo(audio)

	if err == nil {
		log.Println("insert successful")
		log.Println(audioId)
		audio.ID = audioId

		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    audio,
		})
		return nil
	}

	return err
}

func MusicUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	filename, err := musicUpload(c, file)
	if err != nil {
		return err
	}
	userid := c.Param("userid")
	err = music.MusicUpdate(filename, userid)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    filename,
		"message": "音频上传成功：" + filename,
	})
	return nil
}

func musicUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	return fileapi.FileUpload(c, file, config.MUSCI_PATH)
}
