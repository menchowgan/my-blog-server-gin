package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeCalculator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		log.Println("---start time---", start)
		ctx.Next()
		since := time.Since(start)
		log.Printf("---path: %s, cost time: %v", ctx.FullPath(), since)
	}
}
