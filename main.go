package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	person "gmc-blog-server/api/Person"
	photos "gmc-blog-server/api/Photos"
	db "gmc-blog-server/db"
	model "gmc-blog-server/model"
	router "gmc-blog-server/router"
)

func main() {
	r := gin.Default()

	err := db.InitDB()

	if err != nil {
		panic(err)
	}

	initTables()

	defer func() {
		db.DB.DbRClose()
		db.DB.DbWClose()
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	router.Get(r, "/hello", func(ctx *gin.Context) error {
		ctx.JSON(http.StatusOK, gin.H{
			"message": ctx.Request.Header,
		})
		return nil
	})

	groupMap := router.GroupStruct{
		Group: router.GroupMap{
			"/user": {
				{
					Url:     "/person-info-post",
					Method:  http.MethodPost,
					Handler: person.PersonInfoPost,
				}, {
					Url:     "/get-user-simple-info/:id",
					Method:  http.MethodGet,
					Handler: person.GerUserSimpleInfo,
				}, {
					Url:     "/search-user-brief/:id",
					Method:  http.MethodGet,
					Handler: person.GerUserBriefInfo,
				},
			},
			"/photo": {
				{
					Url:     "/avatar/upload/:userid",
					Method:  http.MethodPost,
					Handler: photos.AvatarUpload,
				}, {
					Url:     "/user/photos/upload/:userid", //将用户的图片列表组成字符串存到用户响应表里
					Method:  http.MethodPost,
					Handler: photos.UserPhotosUpload,
				}, {
					Url:     "/user/photos/delete",
					Method:  http.MethodDelete,
					Handler: photos.UserPhotosDelete,
				},
			},
		},
	}

	router.Group(r, groupMap)

	r.NoRoute(func(ctx *gin.Context) {
		fmt.Println("------no route")
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "陆游不存在",
		})
	})

	r.Run(":8888")
}

func initTables() {
	dw := db.DB.GetDbW()
	dr := db.DB.GetDbR()

	var err error
	has := dr.Migrator().HasTable(&model.User{})

	if !has {
		err = dw.AutoMigrate(&model.User{})
	}

	if err != nil {
		panic(err)
	}

	has = dr.Migrator().HasTable(&model.Articles{})

	if !has {
		err = dw.AutoMigrate(&model.Articles{})
	}

	if err != nil {
		panic(err)
	}

	has = dr.Migrator().HasTable(&model.Photos{})

	if !has {
		err = dw.AutoMigrate(&model.Photos{})
	}

	if err != nil {
		panic(err)
	}
}
