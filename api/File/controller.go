package fileapi

import (
	"gmc-blog-server/config"
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
	Ctx      *gin.Context
	Code     int
}

func FileUpload(c *gin.Context, file *multipart.FileHeader, path string) (string, error) {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	folderName := userid
	folderPath := filepath.Join(path, folderName)
	// 使用 0755 权限，更安全
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0755)
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

func Check(c *gin.Context) error {
	// 使用配置文件中的路径
	return DoUpload(
		c,
		config.AppConfig.ChunkPath,
		config.AppConfig.FilePath,
		func(cu *ChunkUpload) {
			if cu.Total == 1 {
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"success": true,
					"data":    "coverVideo/" + cu.FileName,
				})
			}
		})
}

func MergeChunks(c *gin.Context) error {
	// 使用配置文件中的路径
	return DoMerge(
		c,
		config.AppConfig.ChunkPath,
		config.AppConfig.FilePath,
		func(cu *ChunkUpload) {
			response.Success(nil, "File merged successfully", c)
		})
}

func DoUpload(c *gin.Context, basePath, dstPath string, handler func(*ChunkUpload)) error {
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
		BasePath: basePath,
		DstPath:  dstPath,
		Ctx:      c,
	}

	err = chunk.UploadChunk()
	if err != nil {
		return err
	}

	handler(&chunk)

	return err
}

func DoMerge(c *gin.Context, basePath, dstPath string, handler func(*ChunkUpload)) error {
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
		BasePath: basePath,
		DstPath:  dstPath,
		Ctx:      c,
	}

	err = chunk.Merge()
	if err != nil {
		return err
	}

	handler(&chunk)

	return nil
}

func (cu *ChunkUpload) UploadChunk() error {
	// 为每个文件创建一个独立的临时分片目录，防止冲突
	chunkDir := filepath.Join(cu.BasePath, cu.FileName)
	DirNotExistMkdir(chunkDir)
	DirNotExistMkdir(cu.DstPath)

	// 秒传检查：检查最终文件是否存在
	exist := IsExist(filepath.Join(cu.DstPath, cu.FileName))
	if exist {
		cu.Code = response.FileExist
		return nil
	}

	// 分片保存路径
	dst := filepath.Join(chunkDir, cu.Index)

	startFrom := 0
	var err error

	if cu.Total == 1 {
		dst = filepath.Join(cu.DstPath, cu.FileName)
	} else {
		// 断点续传检查：检查分片临时目录中已存在的分片数量
		log.Println("chunk name to save: ", cu.FileName)
		pattern := filepath.Join(chunkDir, "*") // 在文件独立的目录里查找
		startFrom, err = FileExist(pattern)
		if err != nil {
			return err
		}
	}

	log.Println("dst: ", dst)
	// 检查当前分片是否存在
	if IsExist(dst) {
		if cu.Total > 1 {
			cu.Ctx.JSON(http.StatusOK, gin.H{
				"data":      "exist",
				"msg":       dst + " exists",
				"startFrom": startFrom,
			})
		}
		return nil
	}

	file, _ := cu.Ctx.FormFile("file")
	log.Println("file size", file.Size)
	err = cu.Ctx.SaveUploadedFile(file, dst)
	if err != nil && file.Filename != "" {
		return err
	}
	if cu.Total > 1 {
		cu.Ctx.JSON(http.StatusOK, gin.H{
			"data":      "success",
			"msg":       dst,
			"startFrom": startFrom + 1,
		})
	}
	return nil
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func DirNotExistMkdir(path string) {
	// 使用 0755 权限，更安全
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

func (cu *ChunkUpload) Merge() error {
	var err error

	dstFilepath := filepath.Join(cu.DstPath, cu.FileName)
	// 为每个文件创建一个独立的临时分片目录
	chunkDir := filepath.Join(cu.BasePath, cu.FileName)

	// 根据文件名创建文件
	file, err := os.Create(dstFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < cu.Total; i++ {
		// 从文件独立的目录里读取分片
		baseFilepath := filepath.Join(chunkDir, strconv.Itoa(i))
		// 读取分片文件
		chunkFile, err := os.OpenFile(baseFilepath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, chunkFile); err != nil {
			chunkFile.Close()
			return err
		}
		chunkFile.Close()
	}

	// [修复] 合并成功后，删除对应的临时分片目录
	err = os.RemoveAll(chunkDir)
	if err != nil {
		return err
	}

	return nil
}

func FileExist(pattern string) (int, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return 0, err
	}

	return len(files), err
}
