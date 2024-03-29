package music

import (
	fileapi "gmc-blog-server/api/File"
	"gmc-blog-server/config"
	"gmc-blog-server/model"
	"gmc-blog-server/response"
	music "gmc-blog-server/view/Music"
	photos "gmc-blog-server/view/Photos"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Query(c *gin.Context) error {
	userid := c.Param("id")
	mu, err := music.SearchByUserId(userid)

	if err != nil {
		return err
	}

	response.Success(mu, "查询 music 成功", c)
	return nil
}

func MusicCoverUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	log.Println("photo file: ", file.Filename)
	if err == nil {
		filename, err := coverUpload(c, file)
		if err == nil && filename != "" {
			response.Success("cover/"+file.Filename, "", c)
			return nil
		}
		response.Fail(
			photos.AvatarUploadFailed,
			nil,
			photos.StatusText(photos.AvatarUploadFailed),
			c)
		return nil
	}
	return err
}

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

		response.Success(audio, "接受成功", c)
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
	response.Success(filename, "音频上传成功："+filename, c)
	return nil
}

func musicUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	return fileapi.FileUpload(c, file, config.MUSCI_PATH)
}

func coverUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	folderName := "cover"
	folderPath := filepath.Join(config.PHOTO_PATH, userid, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	}
	log.Println(file.Filename)
	dst := filepath.Join(folderPath, file.Filename)

	log.Printf("file path: %s\n", folderPath)
	log.Printf("file name : %s", file.Filename)

	err := c.SaveUploadedFile(file, dst)
	if err == nil && file.Filename != "" {
		return file.Filename, err
	}
	return "", err
}
