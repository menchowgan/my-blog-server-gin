package user

import (
	Article "gmc-blog-server/view/Article"
	photos "gmc-blog-server/view/Photos"

	"gorm.io/gorm"
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
	ArticleSimplaeInfos []Article.ArticleSimpleInfoModel `json:"articleSimplaeInfos"`
}

type UserEnrollModel struct {
	gorm.Model
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Hobbies  string `json:"hobbies"`
	Fans     int    `json:"fans"`
	Evaluate int    `json:"evaluate"`
	Brief    string `json:"brief"`
	Avatar   string `json:"avatar"`
}

type PsersonSimpleIinfo struct {
	ID       uint               `json:"id"`
	Nickname string             `json:"nickname"`
	Fans     int16              `json:"fans"`
	Avatar   string             `json:"avatar"`
	Photos   []photos.PhotoInfo `json:"photos"`
}

type PsersonBrief struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Hobbies  string `json:"hobbies"`
	Avatar   string `json:"avatar"`
}
