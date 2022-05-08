package article

import "time"

type ArticleSimpleInfoModel struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	ImgUrl  string    `json:"imgUrl"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
	Type    string    `json:"type"`
}
