package logic

import (
	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"
	"app/pkg/xerr"
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutInList 获取某个用户好友申请列表
func (l *FriendPutInListLogic) FriendPutInList(in *rpc.FriendPutInListReq) (*rpc.FriendPutInListResp, error) {
	friendRequests, err := l.svcCtx.FriendRequestsModel.ListNoHandler(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "find frends request list err by userId: %s, err: %v", in.UserId, err)
	}

	var resp []*rpc.FriendRequest
	copier.Copy(&resp, friendRequests)
	return &rpc.FriendPutInListResp{
		List: resp,
	}, nil
}
