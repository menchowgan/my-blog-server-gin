package video

import (
	fileapi "gmc-blog-server/api/File"
	"gmc-blog-server/config"
	"gmc-blog-server/model"
	"gmc-blog-server/response"
	photos "gmc-blog-server/view/Photos"
	video "gmc-blog-server/view/Video"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Query(c *gin.Context) error {
	userid := c.Param("id")
	vi, err := video.SearchByUserId(userid)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "查询 video 成功",
		"data":    vi,
	})
	return nil
}

func VideoUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	filename, err := videoUpload(c, file)
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

func VideoUploadChunk(c *gin.Context) error {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	if userid == "" {
		response.Fail(http.StatusBadRequest, nil, "Invalid userid", c)
		return nil
	}

	return fileapi.DoUpload(
		c,
		config.VIDEO_PATH+userid+"/chunks",
		config.VIDEO_PATH+userid,
		func(cu *fileapi.ChunkUpload) {
			if cu.Code == response.FileExist {
				cu.Ctx.JSON(http.StatusOK, gin.H{
					"data": cu.FileName,
					"msg":  cu.FileName + response.StatusText(response.FileExist),
					"code": response.FileExist,
				})
				return
			}
			if cu.Total == 1 {
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"data":    cu.FileName,
					"message": "音频上传成功：" + cu.FileName,
				})
			}
		},
	)
}

func VideoUploadMerge(c *gin.Context) error {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	if userid == "" {
		response.Fail(http.StatusBadRequest, nil, "Invalid userid", c)
		return nil
	}

	return fileapi.DoMerge(
		c,
		config.VIDEO_PATH+userid+"/chunks",
		config.VIDEO_PATH+userid,
		func(cu *fileapi.ChunkUpload) {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"data":    cu.FileName,
				"message": "音频上传成功：" + cu.FileName,
			})
		},
	)
}

func videoUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	return fileapi.FileUpload(c, file, config.VIDEO_PATH)
}

func VideoCoverUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	log.Println("photo file: ", file.Filename)
	if err != nil {
		return err
	}
	filename, err := coverUpload(c, file)
	if err == nil && filename != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
			"data":    "coverVideo/" + file.Filename,
		})
		return nil
	}
	c.JSON(photos.AvatarUploadFailed, gin.H{
		"code":    photos.AvatarUploadFailed,
		"data":    nil,
		"message": photos.StatusText(photos.AvatarUploadFailed),
	})
	return nil
}

func VideoCoverUploadChunk(c *gin.Context) error {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	if userid == "" {
		response.Fail(http.StatusBadRequest, nil, "Invalid userid", c)
		return nil
	}

	return fileapi.DoUpload(
		c,
		config.VIDEO_PATH+userid+"/coverVideo/chunks",
		config.VIDEO_PATH+userid+"/coverVideo",
		func(cu *fileapi.ChunkUpload) {
			if cu.Total == 1 {
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"success": true,
					"data":    "coverVideo/" + cu.FileName,
				})
			}
		},
	)
}

func VideoCoverUploadMerge(c *gin.Context) error {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	if userid == "" {
		response.Fail(http.StatusBadRequest, nil, "Invalid userid", c)
		return nil
	}

	return fileapi.DoMerge(
		c,
		config.VIDEO_PATH+userid+"/coverVideo/chunks",
		config.VIDEO_PATH+userid+"/coverVideo",
		func(cu *fileapi.ChunkUpload) {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"success": true,
				"data":    "coverVideo/" + cu.FileName,
			})
		},
	)
}

func UserVideoUpload(c *gin.Context) error {
	var videoModel model.Video

	if err := c.ShouldBind(&videoModel); err != nil {
		return err
	}
	log.Println("video post: ", videoModel.VideoUrl)
	log.Println("video post: ", videoModel.UserId)
	log.Println("video post: ", videoModel.Avatar)
	log.Println("video post: ", videoModel.Title)
	log.Println("video post: ", videoModel.Artist)
	log.Println("video post: ", videoModel.Evalution)

	err := video.InsertVideoInfo(&videoModel)

	if err == nil {
		log.Println("insert successful")
		log.Println(videoModel)

		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    videoModel,
		})
		return nil
	}

	return err
}

func coverUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	folderName := "coverVideo"
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
