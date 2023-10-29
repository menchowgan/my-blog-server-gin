package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	code    int
	data    interface{}
	message string
	error   bool
}

func Get(code int, data interface{}, msg string, error bool, c *gin.Context) {
	c.JSON(code, gin.H{
		"code":    code,
		"data":    data,
		"message": msg,
		"error":   error,
	})
}

func Success(data interface{}, msg string, c *gin.Context) {
	Get(http.StatusOK, data, msg, true, c)
}

func ServerError(data interface{}, msg string, c *gin.Context) {
	Get(http.StatusInternalServerError, data, msg, false, c)
}

func Fail(code int, data interface{}, msg string, c *gin.Context) {
	Get(code, data, msg, true, c)
}
