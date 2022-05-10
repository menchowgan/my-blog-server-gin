package photos

const (
	AvatarUploadFailed = 900
)

var statusText = map[int]string{
	AvatarUploadFailed: "头像上传失败",
}

func StatusText(code int) string {
	return statusText[code]
}
