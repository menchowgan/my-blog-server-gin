package person

import (
	model "gmc-blog-server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PersonInfoPost(c *gin.Context) error {
	var person model.PersonInfoModel

	if err := c.ShouldBind(&person); err != nil {
		return err
	}

	//todo

	c.JSON(http.StatusOK, gin.H{
		"message": "接受成功",
		"code":    0,
		"data":    "",
	})

	return nil
}
