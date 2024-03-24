package router

import (
	errhandler "gmc-blog-server/errHandler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GroupMap map[string][]struct {
	Url        string
	Method     string
	Handler    Handler
	Middleware gin.HandlerFunc
}

type GroupStruct struct {
	Group GroupMap
}

func Group(r *gin.Engine, groupMap GroupStruct) {
	for key, group := range groupMap.Group {
		routerGroup := r.Group(key)
		{
			for _, route := range group {
				middleware := route.Middleware
				if middleware == nil {
					middleware = func(ctx *gin.Context) {
						ctx.Next()
					}
				}
				switch route.Method {
				case http.MethodGet:
					get(routerGroup, route.Url, middleware, route.Handler)
				case http.MethodPost:
					post(routerGroup, route.Url, middleware, route.Handler)
				case http.MethodPut:
					put(routerGroup, route.Url, middleware, route.Handler)
				case http.MethodDelete:
					delete(routerGroup, route.Url, middleware, route.Handler)
				}
			}
		}
	}
}

func get(rg *gin.RouterGroup, url string, middleware gin.HandlerFunc, callback Handler) {
	rg.GET(url, middleware, func(c *gin.Context) {
		if err := callback(c); err != nil {
			errhandler.Handle(err, c)
		}
	})
}

func post(rg *gin.RouterGroup, url string, middleware gin.HandlerFunc, callback Handler) {
	rg.POST(url, middleware, func(c *gin.Context) {
		if err := callback(c); err != nil {
			errhandler.Handle(err, c)
		}
	})
}

func put(rg *gin.RouterGroup, url string, middleware gin.HandlerFunc, callback Handler) {
	rg.PUT(url, middleware, func(c *gin.Context) {
		if err := callback(c); err != nil {
			errhandler.Handle(err, c)
		}
	})
}

func delete(rg *gin.RouterGroup, url string, middleware gin.HandlerFunc, callback Handler) {
	rg.DELETE(url, middleware, func(c *gin.Context) {
		if err := callback(c); err != nil {
			errhandler.Handle(err, c)
		}
	})
}
