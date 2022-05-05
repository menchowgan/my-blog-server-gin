package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	router "gmc-blog-server/router"

	model "gmc-blog-server/model"
)

func main() {
	fmt.Println("GMC BLOG")
	r := gin.Default()

	router.Get(r, "/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "gmc",
		})
	})

	router.Post(r, "/hello", func(c *gin.Context) {
		var person model.PersonInfoModel

		if err := c.ShouldBind(&person); err != nil {
			c.JSON(400, gin.H{
				"msg": err,
			})
		}
	})

	r.Run(":8888")
}
