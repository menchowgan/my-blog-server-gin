package errhandler

import (
	"fmt"
	"gmc-blog-server/response"
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
		response.Fail(code, nil, err.Error(), c)
	case os.IsPermission(err):
		code = http.StatusForbidden
		response.Fail(code, nil, err.Error(), c)
	case os.IsTimeout(err):
		code = http.StatusRequestTimeout
		response.Fail(code, nil, err.Error(), c)
	default:
		code = http.StatusInternalServerError
		response.Fail(code, nil, err.Error(), c)
	}
}
