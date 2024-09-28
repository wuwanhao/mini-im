package logic

import (
	"app/pkg/xerr"
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendList 获取某个用户的好友列表
func (l *FriendListLogic) FriendList(in *rpc.FriendListReq) (*rpc.FriendListResp, error) {

	friends, err := l.svcCtx.FriendsModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "find frends err by userId: %s, err: %v", in.UserId, err)
	}

	var respList []*rpc.Friends
	copier.Copy(&respList, friends)
	return &rpc.FriendListResp{
		List: respList,
	}, nil
}
