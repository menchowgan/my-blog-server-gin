package router

import (
	article "gmc-blog-server/api/Article"
	fileapi "gmc-blog-server/api/File"
	music "gmc-blog-server/api/Music"
	person "gmc-blog-server/api/Person"
	photos "gmc-blog-server/api/Photos"
	plan "gmc-blog-server/api/Plan"
	video "gmc-blog-server/api/Video"
	weather "gmc-blog-server/api/Weather"
	jwt "gmc-blog-server/token"
	"net/http"
)

func CreateRouter() GroupStruct {
	return GroupStruct{
		Group: GroupMap{
			"/user": {
				{
					Url:     "/person-info-post", // 注册用户信息
					Method:  http.MethodPost,
					Handler: person.PersonInfoPost,
				}, {
					Url:        "/get-user-simple-info", // 首页获取用户简单信息 用于管理端
					Method:     http.MethodGet,
					Handler:    person.GerUserSimpleInfo,
					Middleware: jwt.JwtMiddleware(),
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
					Url:        "/avatar/upload/:userid", // 用户头像上传
					Method:     http.MethodPost,
					Handler:    photos.AvatarUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/user/photos/upload/:userid", // 将用户的图片列表组成字符串存到用户响应表里
					Method:     http.MethodPost,
					Handler:    photos.UserPhotosUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/user/photos/delete", // 先对数据库进行更新再删除文件
					Method:     http.MethodDelete,
					Handler:    photos.UserPhotosDelete,
					Middleware: jwt.JwtMiddleware(),
				},
			},
			"/music": {
				{
					Url:        "/upload/:userid", // 用户收藏音乐上传
					Method:     http.MethodPost,
					Handler:    music.MusicUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/cover/upload/:userid", // 用户上传音乐的封面
					Method:     http.MethodPost,
					Handler:    music.MusicCoverUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/user/upload", // 用户收藏歌曲完整信息上传
					Method:     http.MethodPost,
					Handler:    music.UserMusicUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:     "/query/:id",
					Method:  http.MethodPost,
					Handler: music.Query,
				},
			},
			"/article": {
				{
					Url:        "/avatar/upload/:userid", // 文章封面图片上传
					Method:     http.MethodPost,
					Handler:    article.ArticleAvatarUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/photo/upload/:userid", // 文章内图片上传
					Method:     http.MethodPost,
					Handler:    article.ArticlePhotosUPload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/video/upload/:userid", // 文章内视频上传
					Method:     http.MethodPost,
					Handler:    article.ArticleVideoUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/upload", // 文章整体完整信息上传
					Method:     http.MethodPost,
					Handler:    article.ArticlePost,
					Middleware: jwt.JwtMiddleware(),
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
					Url:        "/upload/:userid",
					Method:     http.MethodPost,
					Handler:    video.VideoUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/cover/upload/:userid", // 用户上传视频的封面
					Method:     http.MethodPost,
					Handler:    video.VideoCoverUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:        "/user/upload",
					Method:     http.MethodPost,
					Handler:    video.UserVideoUpload,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:     "/query/:id",
					Method:  http.MethodPost,
					Handler: video.Query,
				},
			},
			"plan": {
				{
					Url:        "/submit",
					Method:     http.MethodPost,
					Handler:    plan.Submit,
					Middleware: jwt.JwtMiddleware(),
				}, {
					Url:     "/search/:userId",
					Method:  http.MethodGet,
					Handler: plan.Search,
				},
			},
			"weather": {
				{
					Url:     "/query", // 根据城市查询天气
					Method:  http.MethodGet,
					Handler: weather.GetWeather,
				},
			},
			"/file": {
				{
					Url:     "/upload",
					Method:  http.MethodPost,
					Handler: fileapi.Check,
				},
				{
					Url:     "/merge",
					Method:  http.MethodPost,
					Handler: fileapi.MergeChunks,
				},
			},
		},
	}
}
