package model

type PersonInfoModel struct {
	ID                  int64                    `json:"id"`
	Nickname            string                   `json:"nickname"`
	Gender              rune                     `json:"gender"`
	Hobbies             string                   `json:"hobbies"`
	Fans                int                      `json:"fans"`
	Evaluate            int                      `json:"evaluate"`
	Brief               string                   `json:"brief"`
	Avatar              string                   `json:"avatar"`
	Photos              string                   `json:"photos"`
	ArticleSimplaeInfos []ArticleSimpleInfoModel `json:"articleSimplaeInfos"`
}

type User struct {
	ID       int64  `gorm:"column:id;not null;type:int(4) primary key auto_increment;comment:'用户id'"`
	Nickname string `gorm:"column:nickname;type:string;comment:'昵称'"`
	Hobbies  string `gorm:"column:hobbies;type:string;comment:'兴趣'"`
	Gender   rune   `gorm:"column:gender;type:varchar(1);comment:'性别'"`
	Fans     int16  `gorm:"column:fans;type:int(4);comment:'粉丝数量'"`
	Evaluate int16  `gorm:"column:evaluate;type:int(4);comment:'好评'"`
	Brief    string `gorm:"column:brief;type:string;comment:'简介'"`
	Avatar   string `gorm:"column:avatar;type:string;comment:'头像'"`
}
