package person

import (
	"fmt"
	"gmc-blog-server/model"
	user "gmc-blog-server/view/User"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PersonInfoPost(c *gin.Context) error {
	var person model.User

	var err error
	if err = c.ShouldBind(&person); err != nil {
		return err
	}

	if person.ID > 0 {
		person, err := user.Save(&person)

		if err != nil {
			return err
		}

		log.Println("update successful")
		log.Println(person)
		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    person,
		})

		return nil
	} else {
		err := user.InsertUser(&person)

		if err != nil {
			return err
		}
		log.Println("insert successful")
		log.Println(person.ID)

		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    person,
		})

		return nil
	}
}

func GerUserSimpleInfo(c *gin.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	user, err := user.GerUserInfo(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": user,
	})

	return nil
}

func GerUserBriefInfo(c *gin.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	user, err := user.SearchUserBrief(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": user,
	})

	return nil
}
