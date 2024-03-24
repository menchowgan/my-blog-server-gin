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
	ID        int64     `json:"id"`
	UserId    uint      `json:"userId"`
	Title     string    `json:"title"`
	ImgUrl    string    `json:"imgUrl"`
	Brief     string    `json:"brief"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
