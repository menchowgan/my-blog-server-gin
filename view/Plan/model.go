package plan

type PlanModel struct {
	ID      uint   `json:"id"`
	UserId  uint   `json:"userId"`
	Content string `json:"content"`
}
