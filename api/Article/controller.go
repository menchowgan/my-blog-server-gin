package article

import (
	"gmc-blog-server/config"
	"gmc-blog-server/model"
	article "gmc-blog-server/view/Article"
	photos "gmc-blog-server/view/Photos"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Query(c *gin.Context) error {
	userid := c.Param("userid")
	a, err := article.ArticleSimplaeInfosQueryByUserId(userid)
	if err != nil {
		return err
	}
	articleSIs := article.GetArticleInfos(a, userid)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "查询 article 成功",
		"data":    articleSIs,
	})
	return nil
}

func ArticlePost(c *gin.Context) error {
	var a model.Articles

	var err error
	if err = c.ShouldBind(&a); err != nil {
		return err
	}

	if a.ID > 0 {
		err = article.Save(&a)

		if err != nil {
			return err
		}

		log.Println("update successful")
		log.Println(a)
		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    a,
		})

		return nil
	} else {
		err = article.InsertArticle(&a)
		if err != nil {
			return err
		}

		log.Println("insert article successful")
		log.Println(a)
		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    a,
		})
		return nil
	}
}

type articleUploadResponse struct {
	Url  string `json:"url"`
	Alt  string `json:"alt"`
	Href string `json:"href"`
}

func ArticlePhotosUPload(c *gin.Context) error {
	file, err := c.FormFile("wangeditor-uploaded-image")
	if err == nil {
		savePhotoFile(c, file)
	}
	return err
}

func ArticleQuery(c *gin.Context) error {
	aid := c.Param("articleId")
	log.Println("article id: ", aid)

	article, err := article.SearchArticleInfo(aid)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": article,
	})

	return nil
}

func ArticleQueryByType(c *gin.Context) error {
	aType := c.Param("type")
	userid := c.Param("userid")
	log.Println("article type: ", aType)

	article, err := article.SearchArticleInfoByType(userid, aType)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": article,
	})

	return nil
}

func ArticleAvatarUpload(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err == nil {
		savePhotoFile(c, file)
	}
	return err
}

func ArticleVideoUpload(c *gin.Context) error {
	file, err := c.FormFile("wangeditor-uploaded-video")
	if err != nil {
		return err
	}
	saveVideo(c, file)
	return nil
}

func savePhotoFile(c *gin.Context, file *multipart.FileHeader) error {
	log.Println("photo file: ", file.Filename)
	userid := c.Param("userid")
	filename, err := articlePhotoUpload(c, file)
	if err == nil && filename != "" {
		c.JSON(http.StatusOK, gin.H{
			"errno":   0,
			"success": true,
			"data": articleUploadResponse{
				Url:  config.PHOTO_QUERY_PATH + userid + "/article/" + file.Filename,
				Alt:  file.Filename,
				Href: file.Filename,
			},
		})
		return nil
	}
	c.JSON(photos.AvatarUploadFailed, gin.H{
		"code":    photos.AvatarUploadFailed,
		"errno":   1, // 只要不等于 0 就行
		"message": "失败信息",
	})
	return nil
}

func saveVideo(c *gin.Context, file *multipart.FileHeader) error {
	log.Println("video name: ", file.Filename)
	userid := c.Param("userid")
	filename, err := articleVideoUpload(c, file)
	if err != nil || filename == "" {
		c.JSON(photos.AvatarUploadFailed, gin.H{
			"code":    photos.AvatarUploadFailed,
			"errno":   1, // 只要不等于 0 就行
			"message": "失败信息",
		})
		return nil
	}
	c.JSON(http.StatusOK, gin.H{
		"errno":   0,
		"success": true,
		"data": articleUploadResponse{
			Url: config.VIDEO_QUERY_PATH + userid + "/article/" + file.Filename,
		},
	})
	return nil
}

func articlePhotoUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	userid := c.Param("userid")
	log.Println("photo upload user id: ", userid)
	folderName := "article"
	folderPath := filepath.Join(config.PHOTO_PATH, userid, folderName)
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

func articleVideoUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	userid := c.Param("userid")
	folderName := "article"
	folderPath := filepath.Join(config.VIDEO_PATH, userid, folderName)
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
