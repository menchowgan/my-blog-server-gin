package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	article "gmc-blog-server/api/Article"
	music "gmc-blog-server/api/Music"
	person "gmc-blog-server/api/Person"
	photos "gmc-blog-server/api/Photos"
	video "gmc-blog-server/api/Video"
	plan "gmc-blog-server/api/Plan"
	db "gmc-blog-server/db"
	router "gmc-blog-server/router"
)

func main() {
	r := gin.Default()

	err := db.InitDB()

	if err != nil {
		panic(err)
	}

	defer func() {
		db.DB.DbRClose()
		db.DB.DbWClose()
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	db.InitTables()

	router.Get(r, "/hello", func(ctx *gin.Context) error {
		ctx.JSON(http.StatusOK, gin.H{
			"message": ctx.Request.Header,
		})
		return nil
	})

	groupMap := router.GroupStruct{
		Group: router.GroupMap{
			"/user": {
				{
					Url:     "/person-info-post", // 注册用户信息
					Method:  http.MethodPost,
					Handler: person.PersonInfoPost,
				}, {
					Url:     "/get-user-simple-info/:id", // 首页获取用户简单信息
					Method:  http.MethodGet,
					Handler: person.GerUserSimpleInfo,
				}, {
					Url:     "/search-user-brief/:id", // 详情页查询用户详细信息
					Method:  http.MethodGet,
					Handler: person.GerUserBriefInfo,
				}, {
					Url:     "/get/:userid", // 根据用户id获取用户信息
					Method:  http.MethodGet,
					Handler: person.GetInfo,
				}, {
					Url:     "/enroll", // 注册用户信息
					Method:  http.MethodPost,
					Handler: person.Enroll,
				}, {
					Url:     "/login", // 用户登录
					Method:  http.MethodPost,
					Handler: person.Login,
				}, {
					Url:     "/simple-life/:id", // simple-life 网站首页查询
					Method:  http.MethodPost,
					Handler: person.FindSimpleInfo,
				},
			},
			"/photo": {
				{
					Url:     "/avatar/upload/:userid", // 用户头像上传
					Method:  http.MethodPost,
					Handler: photos.AvatarUpload,
				}, {
					Url:     "/user/photos/upload/:userid", // 将用户的图片列表组成字符串存到用户响应表里
					Method:  http.MethodPost,
					Handler: photos.UserPhotosUpload,
				}, {
					Url:     "/user/photos/delete", // 先对数据库进行更新再删除文件
					Method:  http.MethodDelete,
					Handler: photos.UserPhotosDelete,
				},
			},
			"/music": {
				{
					Url:     "/upload/:userid", // 用户收藏音乐上传
					Method:  http.MethodPost,
					Handler: music.MusicUpload,
				}, {
					Url:     "/cover/upload/:userid", // 用户上传音乐的封面
					Method:  http.MethodPost,
					Handler: music.MusicCoverUpload,
				}, {
					Url:     "/user/upload", // 用户收藏歌曲完整信息上传
					Method:  http.MethodPost,
					Handler: music.UserMusicUpload,
				}, {
					Url:     "/query/:id",
					Method:  http.MethodPost,
					Handler: music.Query,
				},
			},
			"/article": {
				{
					Url:     "/avatar/upload/:userid", // 文章封面图片上传
					Method:  http.MethodPost,
					Handler: article.ArticleAvatarUpload,
				}, {
					Url:     "/photo/upload/:userid", // 文章内图片上传
					Method:  http.MethodPost,
					Handler: article.ArticlePhotosUPload,
				}, {
					Url:     "/video/upload/:userid", // 文章内视频上传
					Method:  http.MethodPost,
					Handler: article.ArticleVideoUpload,
				}, {
					Url:     "/upload", // 文章整体完整信息上传
					Method:  http.MethodPost,
					Handler: article.ArticlePost,
				}, {
					Url:     "/query/:articleId", // 文章查询 使用id
					Method:  http.MethodGet,
					Handler: article.ArticleQuery,
				}, {
					Url:     "/query-by-type/:userid/:type", // 文章查询 使用类型名进行模糊查询
					Method:  http.MethodGet,
					Handler: article.ArticleQueryByType,
				}, {
					Url:     "/query-by-userid/:userid", // 文章查询 使用userid
					Method:  http.MethodPost,
					Handler: article.Query,
				},
			},
			"/video": {
				{
					Url:     "/upload/:userid",
					Method:  http.MethodPost,
					Handler: video.VideoUpload,
				}, {
					Url:     "/cover/upload/:userid", // 用户上传视频的封面
					Method:  http.MethodPost,
					Handler: video.VideoCoverUpload,
				}, {
					Url:     "/user/upload",
					Method:  http.MethodPost,
					Handler: video.UserVideoUpload,
				}, {
					Url:     "/query/:id",
					Method:  http.MethodPost,
					Handler: video.Query,
				},
			},
      "plan": {
        {
          Url: "/submit",
          Method: http.MethodPost,
          Handler: plan.Submit,
        }, {
          Url: "/search/:userId",
          Method: http.MethodGet,
          Handler: plan.Search,
        },
      },
		},
	}

	router.Group(r, groupMap)

	r.NoRoute(func(ctx *gin.Context) {
		fmt.Println("------no route")
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "陆游不存在",
		})
	})

	r.Run(":8888")
}
