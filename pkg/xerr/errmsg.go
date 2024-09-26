package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器错误",
	REQUEST_PARAM_ERROR: "请求参数错误",
	DB_ERROR:            "数据库错误",
	TOKEN_INVALID_ERROR: "Token错误",
}

func GetErrMsg(errcode int) string {
	if errMsg, ok := codeText[errcode]; ok {
		return errMsg
	}
	return codeText[SERVER_COMMON_ERROR]
}
