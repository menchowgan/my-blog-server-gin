package person

import (
	"fmt"
	"gmc-blog-server/model"
	user "gmc-blog-server/view/User"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PersonInfoPost(c *gin.Context) error {
	var person model.User

	if err := c.ShouldBind(&person); err != nil {
		return err
	}

	uid, err := user.InsertUser(person)

	if err == nil {
		fmt.Println("insert successful")
		fmt.Println(uid)
		person.ID = uid

		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    person,
		})
	}

	return err
}

func GerUserSimpleInfo(c *gin.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	user := user.GerUserInfo(id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": user,
	})

	return nil
}

func GerUserBriefInfo(c *gin.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	user := user.SearchUserBrief(id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": user,
	})

	return nil
}
