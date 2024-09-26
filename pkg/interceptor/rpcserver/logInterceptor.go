package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogInterceptor 错误日志拦截器
func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, err
	}

	logx.WithContext(ctx).Errorf("[RPC SRV ERR] %v", err)
	casedErr := errors.Cause(err)
	if e, ok := casedErr.(*zerr.CodeMsg); ok {
		err = status.Error(codes.Code(e.Code), e.Msg)
	}
	return resp, err
}
