package resultx

import (
	"context"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Data: data,
		Msg:  "",
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Data: nil,
		Msg:  err,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

//func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
//	return func(ctx context.Context, err error) (int, any) {
//		errcode := xerr.SERVER_COMMON_ERROR
//		errmsg := xerr.GetErrMsg(errcode)
//
//		causeErr := errors.Cause(err)
//		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
//			errcode = e.code
//			errmsg = e.msg
//		} else {
//			if gstatus, ok := status.FromError(causeErr); ok {
//				errcode = int(gstatus.Code())
//				errmsg = gstatus.Message()
//			}
//		}
//
//		// 日志记录
//		logx.WithContext(ctx).Errorf("【%s】 err %v", name, err)
//		return http.StatusBadRequest, Fail(errcode, errmsg)
//	}
//}
