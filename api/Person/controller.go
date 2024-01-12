package person

import (
	"fmt"
	"gmc-blog-server/config"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"gmc-blog-server/response"
	jwt "gmc-blog-server/token"
	article "gmc-blog-server/view/Article"
	music "gmc-blog-server/view/Music"
	photos "gmc-blog-server/view/Photos"
	user "gmc-blog-server/view/User"
	video "gmc-blog-server/view/Video"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type searcher = func(string, *user.PsersonSimpleIinfo, *sync.WaitGroup)

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

		response.Success(person, "接受成功", c)

		return nil
	} else {
		err := user.InsertUser(&person)

		if err != nil {
			return err
		}
		log.Println("insert successful")
		log.Println(person.ID)
		response.Success(person, "接受成功", c)
		return nil
	}
}

func GerUserSimpleInfo(c *gin.Context) error {
	id := c.GetInt("userId")
	fmt.Println(id)

	if id <= 0 {
		response.ServerError(nil, "查询id格式错误，需大于零", c)
		return nil
	}

	idNumb := strconv.Itoa(id)

	user, err := user.GerUserInfo(idNumb)
	if err != nil {
		return err
	}

	response.Success(user, "", c)

	return nil
}

func GerUserBriefInfo(c *gin.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	user, err := user.SearchUserBrief(id)
	if err != nil {
		return err
	}

	response.Success(user, "", c)

	return nil
}

func GetInfo(c *gin.Context) error {
	id := c.GetInt("userId")
	u, err := user.GetUserAllInfo(strconv.Itoa(id))
	if err == nil {
		response.Success(u, "查询用户信息成功", c)
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
		response.ServerError(nil, "用户注册失败（密码存储错误）", c)

		return nil
	}

	user := &model.User{}
	user.Nickname = userEnroll.UserName
	user.Password = string(hash)
	log.Println("new user:", user.Nickname, user.Password)

	dw := db.DB.GetDbW()
	transaction := dw.Begin()
	err := transaction.Create(user).Error
	if err != nil {
		transaction.Rollback()
		return err
	}
	response.Success(user.ID, "用户注册成功", c)
	transaction.Commit()
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
		token, err := jwt.CreateToken(int(ur.ID), ur.Nickname)
		if err != nil {
			response.Fail(
				response.TokenCreateFailed,
				nil,
				response.StatusText(response.TokenCreateFailed),
				c,
			)
			return err
		}
		response.Success(token, "用户注册成功", c)
	} else {
		response.Success(nil, "未找到用户信息", c)
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

	if info.ID == 0 {
		response.Get(10010, nil, "用户信息未找到", false, c)
	}

	response.Success(info, "查询simple life信息成功", c)

	return nil
}

func getSimpleLifeInfo(id string) (user.PsersonSimpleIinfo, error) {
	userInfo := user.PsersonSimpleIinfo{}

	fns := []searcher{
		searchPhoto,
		searchMusic,
		searchArticle,
		searchUser,
		searchVideo,
	}

	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)
		go fn(id, &userInfo, &wg)
	}

	wg.Wait()
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

func searchMusic(id string, userInfo *user.PsersonSimpleIinfo, wg *sync.WaitGroup) {
	defer wg.Done()
	musics, _ := music.MusicQueryByUserIdSimplaeLife(id)
	if len(musics) == 0 {
		return
	}
	log.Println("simple info: found music:", musics)
	musicModel := []music.MusicInfo{}
	for _, m := range musics {
		if m.ID != 0 {
			musicModel = append(musicModel, music.MusicInfo{
				ID:        m.ID,
				Artist:    m.Artist,
				Title:     m.Title,
				AudioUrl:  config.MUSCI_QUERY_PATH + id + "/" + m.AudioUrl,
				Evalution: m.Evalution,
				Avatar:    config.MUSCI_QUERY_PATH + id + "/" + m.Avatar,
			})
		}
	}

	userInfo.Musics = musicModel
}

func searchArticle(id string, userInfo *user.PsersonSimpleIinfo, wg *sync.WaitGroup) {
	defer wg.Done()
	oneArticle, _ := article.ArticleQueryByUserIdSimplaeLife(id)
	if oneArticle.ID == 0 {
		return
	}
	log.Println("simple info: found article:", oneArticle)
	oneArticleModel := article.ArticleSimpleInfoModel{
		ID:     int64(oneArticle.ID),
		Brief:  oneArticle.Brief,
		Title:  oneArticle.Title,
		UserId: oneArticle.UserId,
		ImgUrl: config.PHOTO_QUERY_PATH + id + "/" + "article" + "/" + oneArticle.ImgUrl,
		Date:   oneArticle.CreatedAt,
	}

	userInfo.Article = oneArticleModel
}

func searchPhoto(id string, userInfo *user.PsersonSimpleIinfo, wg *sync.WaitGroup) {
	defer wg.Done()
	photosInfo, _ := photos.PhotosQueryByUserId(id)
	if len(photosInfo.ImgUrls) == 0 {
		return
	}
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
}

func searchUser(id string, userInfo *user.PsersonSimpleIinfo, wg *sync.WaitGroup) {
	defer wg.Done()
	userI, _ := user.FindUser(id)
	if userI.ID == 0 {
		return
	}
	userInfo.Avatar = config.PHOTO_QUERY_PATH + id + "/" + userI.Avatar
	userInfo.Nickname = userI.Nickname
	userInfo.ID = userI.ID
}

func searchVideo(id string, userInfo *user.PsersonSimpleIinfo, wg *sync.WaitGroup) {
	defer wg.Done()
	videoI, _ := video.SearchSimpleLifeByUserId(id)
	if videoI.ID == 0 {
		return
	}
	userInfo.Video = videoI
}
