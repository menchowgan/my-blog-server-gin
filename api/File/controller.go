package fileapi

import (
	"gmc-blog-server/response"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChunkUpload struct {
	FileName string
	Index    string
	Total    int
	BasePath string // chunks存储路径
	DstPath  string // 文件存储路径
	ctx      *gin.Context
}

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
	if err != nil && file.Filename != "" {
		return file.Filename, err
	}
	return "", err
}

func Check(c *gin.Context) error {
	fileName := c.PostForm("fileName")
	index := c.PostForm("index")
	totalChunks := c.PostForm("totalChunks")
	total, err := strconv.Atoi(totalChunks)
	if err != nil {
		response.Fail(http.StatusBadRequest, nil, "Invalid totalChunks", c)
		return nil
	}

	chunk := ChunkUpload{
		FileName: fileName,
		Index:    index,
		Total:    total,
		BasePath: "/usr/local/share/chunks",
		DstPath:  "/usr/local/share/chunks",
		ctx:      c,
	}

	return chunk.UploadChunk()
}

func MergeChunks(c *gin.Context) error {
	filename := c.PostForm("fileName")
	if filename == "" {
		response.Fail(http.StatusBadRequest, nil, "Invalid filename", c)
		return nil
	}
	totalChunks := c.PostForm("totalChunks")

	total, err := strconv.Atoi(totalChunks)
	if err != nil {
		response.Fail(http.StatusBadRequest, nil, "Invalid totalChunks", c)
		return nil
	}

	chunk := ChunkUpload{
		FileName: filename,
		Total:    total,
		BasePath: "/usr/local/share/chunks",
		DstPath:  "/usr/local/share/chunks",
		ctx:      c,
	}

	return chunk.Merge()
}

func (cu *ChunkUpload) UploadChunk() error {
	dst := filepath.Join(cu.BasePath, cu.FileName+"-"+cu.Index)

	startFrom := 0
	var err error

	if cu.Total == 1 {
		dst = filepath.Join(cu.DstPath, cu.FileName)
	} else {
		log.Println("chunk name to save: ", cu.FileName)
		pattern := filepath.Join(cu.BasePath, cu.FileName+"*")
		startFrom, err = FileExist(pattern)
		if err != nil {
			return err
		}
	}

	log.Println("dst: ", dst)
	// 检查文件是否存在
	if IsExist(dst) {
		cu.ctx.JSON(http.StatusOK, gin.H{
			"data":      "exist",
			"msg":       dst + " exists",
			"startFrom": startFrom,
		})
		return nil
	}

	file, _ := cu.ctx.FormFile("file")
	log.Println("file size", file.Size)
	err = cu.ctx.SaveUploadedFile(file, dst)
	if err != nil && file.Filename != "" {
		return err
	}
	cu.ctx.JSON(http.StatusOK, gin.H{
		"data":      "success",
		"msg":       dst,
		"startFrom": startFrom + 1,
	})
	return nil
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (cu *ChunkUpload) Merge() error {
	var err error

	dstFilepath := filepath.Join(cu.DstPath, cu.FileName)

	// 根据文件名创建文件
	file, err := os.Create(dstFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < cu.Total; i++ {
		baseFilepath := filepath.Join(cu.BasePath, cu.FileName+"-"+strconv.Itoa(i))
		// 读取分片文件
		chunkFile, err := os.Open(baseFilepath)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		if _, err := io.Copy(file, chunkFile); err != nil {
			return err
		}

		// 删除分片文件
		err = os.Remove(baseFilepath)
		if err != nil {
			return err
		}
	}

	response.Success(nil, "File merged successfully", cu.ctx)
	return nil
}

func FileExist(pattern string) (int, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return 0, err
	}

	return len(files), err
}
