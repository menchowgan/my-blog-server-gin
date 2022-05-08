package photos

import (
	"gmc-blog-server/common/ResponseCode"
	"gmc-blog-server/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AvatarUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err == nil {
		log.Println(file.Filename)
		dst := config.PHOTO_PATH + file.Filename

		c.SaveUploadedFile(file, dst)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
			"data":    file.Filename,
		})
		return nil
	}
	c.JSON(ResponseCode.AvatarUpload.Code, gin.H{
		"code":    ResponseCode.AvatarUpload.Code,
		"data":    nil,
		"message": ResponseCode.AvatarUpload.Message,
	})
	return nil
}
