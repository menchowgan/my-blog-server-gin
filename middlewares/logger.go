package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LogValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Loging checker")
		log.Println("path", ctx.FullPath())
		ctx.Next()
	}
}
