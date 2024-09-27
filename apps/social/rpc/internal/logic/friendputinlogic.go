package logic

import (
	"context"

	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 好友业务：请求添加好友、通过或拒绝申请、好有列表
func (l *FriendPutInLogic) FriendPutIn(in *rpc.FriendPutInReq) (*rpc.FriendPutInResp, error) {
	// todo: add your logic here and delete this line

	return &rpc.FriendPutInResp{}, nil
}
