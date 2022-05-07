package errhandler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Handle(err error, c *gin.Context) {
	fmt.Println("----------error", err)
	code := http.StatusOK
	switch {
	case os.IsNotExist(err):
		code = http.StatusNotFound
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
	case os.IsPermission(err):
		code = http.StatusForbidden
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
	default:
		code = http.StatusInternalServerError
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
	}
}
