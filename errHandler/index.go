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
			"message": "没找到该访问方式",
		})
	case os.IsPermission(err):
		code = http.StatusForbidden
		c.JSON(code, gin.H{
			"message": "您没有权限",
		})
	default:
		code = http.StatusInternalServerError
		c.JSON(code, gin.H{
			"message": "系统内部错误",
		})
	}
}
