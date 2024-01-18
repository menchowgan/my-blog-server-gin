package photos

import (
	"encoding/json"
	fileapi "gmc-blog-server/api/File"
	"gmc-blog-server/config"
	"gmc-blog-server/response"
	photos "gmc-blog-server/view/Photos"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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
	body, _ := io.ReadAll(c.Request.Body)

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
			response.Fail(http.StatusNotFound, nil, "文件未找到，无法删除", c)
			return nil
		}
		return err
	}

	response.Success(nil, "删除成功", c)

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

	response.Success(filename, "图片上传成功", c)

	return nil
}

func PhotoUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	log.Println("photo file: ", file.Filename)
	if err == nil {
		filename, err := photoUpload(c, file)
		if err == nil && filename != "" {
			response.Success(file.Filename, "照片上传成功", c)
			return nil
		}
		response.Fail(
			photos.AvatarUploadFailed,
			nil,
			photos.StatusText(photos.AvatarUploadFailed), c)
		return nil
	}
	return err
}

func photoUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	return fileapi.FileUpload(c, file, config.PHOTO_PATH)
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
