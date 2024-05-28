package response

var (
	TokenCreateFailed = -999
	UserInfoNotFound  = 10010
	FileExist         = 10001 // 文件已存在
)

var statusText = map[int]string{
	UserInfoNotFound:  "用户信息未找到",
	TokenCreateFailed: "Token生成失败",
	FileExist:         "文件已上传",
}

func StatusText(code int) string {
	return statusText[code]
}
