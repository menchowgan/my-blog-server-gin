package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	person "gmc-blog-server/api/Person"
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

	router.Get(r, "/hello", func(ctx *gin.Context) error {
		ctx.JSON(http.StatusOK, gin.H{
			"message": ctx.Request.Header,
		})
		return nil
	})

	groupMap := router.GroupStruct{
		Group: router.GroupMap{
			"/user": {{
				Url:     "/person-info-post",
				Method:  http.MethodPost,
				Handler: person.PersonInfoPost,
			}}},
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

	defer func() {
		db.DB.DbRClose()
		db.DB.DbWClose()
	}()

	has := dr.Migrator().HasTable(&model.User{})

	if has {
		return
	}

	err := dw.AutoMigrate(&model.User{})
	if err == nil {
		return
	}
	panic(err)
}
