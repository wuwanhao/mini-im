package logic

import (
	"app/apps/user/rpc/internal/svc"
	"app/apps/user/rpc/user"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserRequest) (*user.FindUserResponse, error) {
	// todo: add your logic here and delete this line

	return &user.FindUserResponse{}, nil
}
