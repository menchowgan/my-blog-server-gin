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
	"strconv"
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

	i := c.GetInt("userId")
	if i == 0 {
		response.Fail(http.StatusForbidden, nil, "重新登录", c)
		return nil
	}
	photos.GetByUserId(strconv.Itoa(i))
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
	photos.GetByUserId(userid)
	response.Success(filename, "图片上传成功", c)

	return nil
}

func UserPhotosUploadChunk(c *gin.Context) error {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	if userid == "" {
		response.Fail(http.StatusBadRequest, nil, "Invalid userid", c)
		return nil
	}

	return fileapi.DoUpload(
		c,
		config.PHOTO_PATH+userid+"/chunk",
		config.PHOTO_PATH+userid,
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

func UserPhotosUploadChunkMerge(c *gin.Context) error {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	if userid == "" {
		response.Fail(http.StatusBadRequest, nil, "Invalid userid", c)
		return nil
	}

	return fileapi.DoMerge(
		c,
		config.PHOTO_PATH+userid+"/chunk",
		config.PHOTO_PATH+userid,
		func(cu *fileapi.ChunkUpload) {
			err := photos.PhotoUpdate(cu.FileName, userid)
			if err != nil {
				response.Fail(http.StatusInternalServerError, nil, "图片上传失败", c)
				return
			}
			photos.GetByUserId(userid)
			response.Success(cu.FileName, "图片上传成功", c)
		},
	)
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
