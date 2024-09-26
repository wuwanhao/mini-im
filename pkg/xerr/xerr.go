package xerr

import "github.com/zeromicro/x/errors"

func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewMsgErr(msg string) error {
	return errors.New(SERVER_COMMON_ERROR, msg)
}

func NewCodeErr(code int) error {
	return errors.New(code, GetErrMsg(code))
}
func NewInternalError() error {
	return errors.New(SERVER_COMMON_ERROR, GetErrMsg(SERVER_COMMON_ERROR))
}

func NewDBError() error {
	return errors.New(DB_ERROR, GetErrMsg(DB_ERROR))
}

func NewReqParamsError() error {
	return errors.New(REQUEST_PARAM_ERROR, GetErrMsg(REQUEST_PARAM_ERROR))
}
