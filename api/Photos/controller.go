package photos

import (
	"encoding/json"
	"gmc-blog-server/config"
	photos "gmc-blog-server/view/Photos"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type deleteUrlBody struct {
	Uri string `json:"url"`
}

func AvatarUpload(c *gin.Context) error {
	return PhotoUpload(c)
}

func UserPhotosDelete(c *gin.Context) error {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var urlBody deleteUrlBody
	err := json.Unmarshal(body, &urlBody)
	uri := urlBody.Uri
	log.Println("url", urlBody.Uri)
	if err != nil {
		return err
	}

	params := strings.Split(uri, "/")
	folderPath, fileName := params[0], params[1]
	uid := folderPath

	err = photos.PhotoDeleteByFileName(uid, fileName)
	if err != nil {
		return err
	}

	err = deletePicture(uri)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "文件未找到，无法删除",
			})
			return nil
		}
		return err
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "删除成功",
	})

	return nil
}

func UserPhotosUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	filename, err := photoUpload(c, file)
	if err != nil {
		return err
	}
	userid := c.Param("userid")
	err = photos.PhotoUpdate(filename, userid)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "图片上传成功",
	})

	return nil
}

func PhotoUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err == nil {
		filename, err := photoUpload(c, file)
		if err == nil && filename != "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"success": true,
				"data":    file.Filename,
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
	return err
}

func photoUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	userid := c.Param("userid")
	folderName := userid
	folderPath := filepath.Join(config.PHOTO_PATH, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	}
	log.Println(file.Filename)
	dst := folderPath + "/" + file.Filename

	log.Printf("file path: %s\n", folderPath)
	log.Printf("file name : %s", file.Filename)

	err := c.SaveUploadedFile(file, dst)
	if err == nil && file.Filename != "" {
		return file.Filename, err
	}
	return "", err
}

func deletePicture(uri string) error {
	filePath := config.PHOTO_PATH + uri
	log.Println("file path: ", filePath)
	if _, err := os.Stat(filePath); err != nil {
		return err
	}
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
