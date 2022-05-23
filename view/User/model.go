package user

import (
	article "gmc-blog-server/view/Article"
	music "gmc-blog-server/view/Music"
	photos "gmc-blog-server/view/Photos"
	video "gmc-blog-server/view/Video"
	"time"
)

type PersonInfoModel struct {
	ID                  uint                             `json:"id"`
	Nickname            string                           `json:"nickname"`
	Gender              string                           `json:"gender"`
	Hobbies             string                           `json:"hobbies"`
	Fans                int16                            `json:"fans"`
	Evaluate            int                              `json:"evaluate"`
	Brief               string                           `json:"brief"`
	Avatar              string                           `json:"avatar"`
	Photos              string                           `json:"photos"`
	Audios              []music.MusicInfo                `json:"audios"`
	Videos              []video.VideoInfo                `json:"videos"`
	ArticleSimplaeInfos []article.ArticleSimpleInfoModel `json:"articleSimplaeInfos"`
}

type UserEnrollModel struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Nickname  string    `json:"nickname"`
	Gender    string    `json:"gender"`
	Hobbies   string    `json:"hobbies"`
	Fans      int       `json:"fans"`
	Evaluate  int       `json:"evaluate"`
	Brief     string    `json:"brief"`
	Avatar    string    `json:"avatar"`
}

type PsersonSimpleIinfo struct {
	ID                 uint                             `json:"id"`
	Nickname           string                           `json:"nickname"`
	Fans               int16                            `json:"fans"`
	Avatar             string                           `json:"avatar"`
	Photos             []photos.PhotoInfo               `json:"photos"`
	ArticleSimpleInfos []article.ArticleSimpleInfoModel `json:"articleSimpleInfos"`
}

type PsersonBrief struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Hobbies  string `json:"hobbies"`
	Avatar   string `json:"avatar"`
}
