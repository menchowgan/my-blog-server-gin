package fileapi

import (
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func FileUpload(c *gin.Context, file *multipart.FileHeader, path string) (string, error) {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	folderName := userid
	folderPath := filepath.Join(path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	}
	log.Println(file.Filename)
	dst := filepath.Join(folderPath, file.Filename)

	log.Printf("file path: %s\n", folderPath)
	log.Printf("file name : %s", file.Filename)

	err := c.SaveUploadedFile(file, dst)
	if err == nil && file.Filename != "" {
		return file.Filename, err
	}
	return "", err
}
