package router

import (
	errhandler "gmc-blog-server/errHandler"

	"github.com/gin-gonic/gin"
)

type GroupMap map[string]struct {
	Url     string
	Method  string
	Handler Handler
}

type GroupStruct struct {
	Group GroupMap
}

func Group(r *gin.Engine, groupMap GroupStruct) {
	routerGroup := r.Group("/user")
	{
		for _, v := range groupMap.Group {
			switch v.Method {
			case "get":
				get(routerGroup, v.Url, v.Handler)
			case "post":
				post(routerGroup, v.Url, v.Handler)
			case "put":
				put(routerGroup, v.Url, v.Handler)
			}
		}
	}
}

func get(rg *gin.RouterGroup, url string, callback Handler) {
	rg.GET(url, func(c *gin.Context) {
		if err := callback(c); err != nil {
			errhandler.Handle(err, c)
		}
	})
}

func post(rg *gin.RouterGroup, url string, callback Handler) {
	rg.POST(url, func(c *gin.Context) {
		if err := callback(c); err != nil {
			errhandler.Handle(err, c)
		}
	})
}

func put(rg *gin.RouterGroup, url string, callback Handler) {
	rg.PUT(url, func(c *gin.Context) {
		if err := callback(c); err != nil {
			errhandler.Handle(err, c)
		}
	})
}
