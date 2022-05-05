package router

import (
	"github.com/gin-gonic/gin"
)

func Get(r *gin.Engine, url string, callback func(*gin.Context)) {
	r.GET(url, callback)
}

func Post(r *gin.Engine, url string, callback func(*gin.Context)) {
	r.POST(url, callback)
}
