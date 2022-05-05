package model

import (
	"time"
)

type ArticleSimpleInfoModel struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	ImgUrl  string    `json:"imgUrl"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
	Type    string    `json:"type"`
}

// interface ArticleSimpleInfoModel {
//   id?: number
//   imgUrl?: string
//   title?: string
//   content?: Date | string
//   date?: Date,
//   type?: string
// }
