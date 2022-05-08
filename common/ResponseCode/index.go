package ResponseCode

type ResponseCode struct {
	Message string
	Code    int
}

var AvatarUpload = ResponseCode{
	Message: "头像上传失败",
	Code:    900,
}
