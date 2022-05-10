package user

const (
	UserInfoUploadFailed = 800
)

var statusText = map[int]string{
	UserInfoUploadFailed: "用户基本信息上传失败",
}

func StatusText(code int) string {
	return statusText[code]
}
