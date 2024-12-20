// Code generated by goctl. DO NOT EDIT.
// Source: im.proto

package server

import (
	"context"

	"app/apps/im/rpc/im"
	"app/apps/im/rpc/internal/logic"
	"app/apps/im/rpc/internal/svc"
)

type ImServer struct {
	svcCtx *svc.ServiceContext
	im.UnimplementedImServer
}

func NewImServer(svcCtx *svc.ServiceContext) *ImServer {
	return &ImServer{
		svcCtx: svcCtx,
	}
}

func (s *ImServer) Ping(ctx context.Context, in *im.Request) (*im.Response, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}
