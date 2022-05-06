package router

import (
	errhandler "gmc-blog-server/errHandler"

	"github.com/gin-gonic/gin"
)

type GroupMap map[string][]struct {
	Url     string
	Method  string
	Handler Handler
}

type GroupStruct struct {
	Group GroupMap
}

func Group(r *gin.Engine, groupMap GroupStruct) {
	for key, group := range groupMap.Group {
		routerGroup := r.Group(key)
		{
			for _, route := range group {
				switch route.Method {
				case "get":
					get(routerGroup, route.Url, route.Handler)
				case "post":
					post(routerGroup, route.Url, route.Handler)
				case "put":
					put(routerGroup, route.Url, route.Handler)
				}
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
