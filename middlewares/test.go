package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func TestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := ctx.Request.Header
		log.Println("cookies: ", ctx.Request.Cookies())
		for key, value := range h {
			log.Println(key, value)
		}
		ctx.Next()
	}
}
