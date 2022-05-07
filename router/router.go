package router

import (
	"fmt"

	"github.com/gin-gonic/gin"

	errhandler "gmc-blog-server/errHandler"
)

type Handler func(*gin.Context) error

func Get(r *gin.Engine, url string, callback Handler) {
	r.GET(url, func(c *gin.Context) {
		if err := callback(c); err != nil {
			fmt.Println("get error", err.Error())
			errhandler.Handle(err, c)
		}
	})
}

func Post(r *gin.Engine, url string, callback Handler) {
	r.POST(url, func(c *gin.Context) {
		if err := callback(c); err != nil {
			fmt.Println("get error", err.Error())
			errhandler.Handle(err, c)
		}
	})
}
