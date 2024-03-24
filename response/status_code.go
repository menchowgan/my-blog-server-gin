package response

var (
	TokenCreateFailed = -999
	UserInfoNotFound  = 10010
)

var statusText = map[int]string{
	UserInfoNotFound:  "用户信息未找到",
	TokenCreateFailed: "Token生成失败",
}

func StatusText(code int) string {
	return statusText[code]
}
