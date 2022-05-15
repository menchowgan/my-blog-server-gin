package article

import (
	"time"
)

type ArticleSimpleInfoModel struct {
	ID     int64     `json:"id"`
	Title  string    `json:"title"`
	UserId uint      `json:"userId"`
	ImgUrl string    `json:"imgUrl"`
	Brief  string    `json:"brief"`
	Date   time.Time `json:"date"`
}

type ArticleInfoModel struct {
	ID      int64     `json:"id"`
	UserId  uint      `json:"userId"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}
