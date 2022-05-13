package music

type MusicInfo struct {
	ID        uint   `json:"id"`
	UserId    int    `json:"userId"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Evalution string `json:"evalution"`
	AudioUrl  string `json:"audioUrl"`
	Avatar    string `json:"avatar"`
}
