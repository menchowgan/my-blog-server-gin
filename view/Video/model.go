package video

type VideoInfo struct {
	ID        uint   `json:"id"`
	UserId    int    `json:"userId"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Evalution string `json:"evalution"`
	VideoUrl  string `json:"videoUrl"`
	Avatar    string `json:"avatar"`
}
