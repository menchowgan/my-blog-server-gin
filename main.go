package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	db "gmc-blog-server/db"
	redis "gmc-blog-server/redis"
	router "gmc-blog-server/router"
  middlewares "gmc-blog-server/middlewares"
)

func main() {
	r := gin.Default()

	err := db.InitDB()

	if err != nil {
		panic(err)
	}

	defer func() {
		db.DB.DbRClose()
		db.DB.DbWClose()
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	db.InitTables()
  redis.Init()

  r.Use(middlewares.LogValidator())
  r.Use(middlewares.TestMiddleware())

	router.Get(r, "/hello", func(ctx *gin.Context) error {
		ctx.JSON(http.StatusOK, gin.H{
			"message": ctx.Request.Header,
		})
		return nil
	})

	groupMap := router.CreateRouter()

	router.Group(r, groupMap)

	r.NoRoute(func(ctx *gin.Context) {
		fmt.Println("------no route")
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "陆游不存在",
		})
	})

	r.Run(":8888")
}
