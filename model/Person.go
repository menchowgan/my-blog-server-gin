package model

type PersonInfoModel struct {
	ID                  int64                    `json:"id"`
	Nickname            string                   `json:"nickname"`
	Gender              string                   `json:"gender"`
	Hobbies             []string                 `json:"hobbies"`
	Fans                int                      `json:"fans"`
	Evaluate            int                      `json:"evaluate"`
	Brief               string                   `json:"brief"`
	Avatar              string                   `json:"avatar"`
	Photos              []string                 `json:"photos"`
	ArticleSimplaeInfos []ArticleSimpleInfoModel `json:"articleSimplaeInfos"`
}
