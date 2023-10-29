package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "gmc-blog-server/db"
	middlewares "gmc-blog-server/middlewares"
	redis "gmc-blog-server/redis"
	"gmc-blog-server/response"
	router "gmc-blog-server/router"
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
		response.Success(nil, "Hi~~~~~~", ctx)
		return nil
	})

	groupMap := router.CreateRouter()

	router.Group(r, groupMap)
	r.NoRoute(func(ctx *gin.Context) {
		response.Fail(http.StatusNotFound, nil, "陆游不存在", ctx)
	})

	r.Run(":8888")
}
