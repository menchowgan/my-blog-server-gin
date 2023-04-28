package person

import (
	"fmt"
	"gmc-blog-server/config"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	article "gmc-blog-server/view/Article"
	music "gmc-blog-server/view/Music"
	photos "gmc-blog-server/view/Photos"
	user "gmc-blog-server/view/User"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PersonInfoPost(c *gin.Context) error {
	var person model.User

	var err error
	if err = c.ShouldBind(&person); err != nil {
		return err
	}

	if person.ID > 0 {
		person, err := user.Save(&person)

		if err != nil {
			return err
		}

		log.Println("update successful")
		log.Println(person)
		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    person,
		})

		return nil
	} else {
		err := user.InsertUser(&person)

		if err != nil {
			return err
		}
		log.Println("insert successful")
		log.Println(person.ID)

		c.JSON(http.StatusOK, gin.H{
			"message": "接受成功",
			"code":    0,
			"data":    person,
		})

		return nil
	}
}

func GerUserSimpleInfo(c *gin.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	user, err := user.GerUserInfo(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": user,
	})

	return nil
}

func GerUserBriefInfo(c *gin.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	user, err := user.SearchUserBrief(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": user,
	})

	return nil
}

func GetInfo(c *gin.Context) error {
	id := c.Param("userid")
	u, err := user.GetUserAllInfo(id)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "查询用户信息成功",
			"data":    u,
		})
	}
	return err
}

func Enroll(c *gin.Context) error {
	var userEnroll model.UserEnroll

	if err := c.ShouldBind(&userEnroll); err != nil {
		return err
	}
	log.Printf("user name: %s; password: %s to enroll! ", userEnroll.UserName, userEnroll.Passwrod)

	hash := createBcryptPassword(userEnroll.Passwrod)
	if hash == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "用户注册失败（密码存储错误）",
		})
		return nil
	}

	user := &model.User{}
	user.Nickname = userEnroll.UserName
	user.Password = string(hash)
	log.Println("new user:", user.Nickname, user.Password)

	dw := db.DB.GetDbW()
	err := dw.Create(user).Error
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "用户注册成功",
		"data":    user.ID,
	})
	return nil
}

func Login(c *gin.Context) error {
	var userLog model.UserEnroll

	if err := c.ShouldBind(&userLog); err != nil {
		return err
	}
	log.Printf("user name: %s; password: %s to login! ", userLog.UserName, userLog.Passwrod)

	u, err := findUser(&userLog)
	if err != nil {
		return err
	}

	var ur model.User

	for _, user := range u {
		log.Printf("----user' s id is %d, name is %s, password is %s\n", user.ID, user.Nickname, user.Password)
		if user.Password == "" {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLog.Passwrod)); err == nil {
			ur = user
		}
	}

	log.Printf("found user' s id is %d name is %s, password is %s, avatar is %s\n", ur.ID, ur.Nickname, ur.Password, ur.Avatar)
	if ur.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "用户登录成功",
			"data":    ur.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "未找到用户信息",
			"data":    nil,
		})
	}
	return nil
}

func FindSimpleInfo(c *gin.Context) error {
	id := c.Param("id")
	log.Println("user id", id)

	info, err := getSimpleLifeInfo(id)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "查询simple life信息成功",
		"data":    info,
	})

	return nil
}

func getSimpleLifeInfo(id string) (user.PsersonSimpleIinfo, error) {
	photosChan := make(chan model.Photos)
	photoError := false

	musicChan := make(chan model.Music)
	musicError := false

	articleChan := make(chan model.Articles)
	articleError := false

	userChan := make(chan model.User)
	userError := false

	go searchPhoto(id, photosChan)
	go searchMusic(id, musicChan)
	go searchArticle(id, articleChan)
	go searchUser(id, userChan)

	userInfo := user.PsersonSimpleIinfo{}

	for {
		if photoError && musicError && articleError && userError {
			break
		}

		select {
		case photosInfo := <-photosChan:
			imgUrls := []photos.PhotoInfo{}
			if photosInfo.ImgUrls != "" && len(photosInfo.ImgUrls) > 0 {
				urls := strings.Split(photosInfo.ImgUrls, ";")
				for idx, url := range urls {
					model := photos.PhotoInfo{
						ID:  uint(idx + 1),
						Url: config.PHOTO_QUERY_PATH + id + "/" + url,
					}
					imgUrls = append(imgUrls, model)
				}
			}

			userInfo.Photos = imgUrls
			photoError = true
		case musics := <-musicChan:
			userInfo.Musics = musics
			musicError = true
		case oneArticle := <-articleChan:
			userInfo.Article = oneArticle
			articleError = true
		case userI := <-userChan:
			userInfo.Avatar = userI.Avatar
			userInfo.Nickname = userI.Nickname
			userError = true
		}
	}

	return userInfo, nil
}

func findUser(user *model.UserEnroll) ([]model.User, error) {
	dr := db.DB.GetDbR()
	var u []model.User

	err := dr.Select("id, nickname, password, avatar").Where("nickname = ?", user.UserName).Find(&u).Error
	if err != nil {
		return []model.User{}, err
	}

	return u, nil
}

func createBcryptPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	log.Println("password hash:", string(hash))
	return string(hash)
}

func searchMusic(id string, musicChan chan model.Music) {
	info, _ := music.MusicQueryByUserIdSimplaeLife(id)
	musicChan <- info
}

func searchArticle(id string, articleChan chan model.Articles) {
	info, _ := article.ArticleQueryByUserIdSimplaeLife(id)
	articleChan <- info
}

func searchPhoto(id string, photoChan chan model.Photos) {
	info, _ := photos.PhotosQueryByUserId(id)
	photoChan <- info
}

func searchUser(id string, userChan chan model.User) {
	info, _ := user.FindUser(id)
	userChan <- info
}
