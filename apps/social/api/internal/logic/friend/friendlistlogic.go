package friend

import (
	"app/apps/social/rpc/socialclient"
	"app/apps/user/rpc/userclient"
	"app/pkg/ctxdata"
	"context"

	"app/apps/social/api/internal/svc"
	"app/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendList 获取某个用户的好友列表
func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// 获取当前登录的用户信息
	uid := ctxdata.GetUid(l.ctx)
	friends, err := l.svcCtx.SocialRpc.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}
	if len(friends.List) == 0 {
		return &types.FriendListResp{}, nil
	}

	// 获取所有的好友 ID
	uids := make([]string, 0, len(friends.List))
	for _, f := range friends.List {
		uids = append(uids, f.FriendUid)
	}

	// 根据 好友 ID 获取所有的好友信息
	user, err := l.svcCtx.UserRpc.FindUser(l.ctx, &userclient.FindUserRequest{
		Ids: uids,
	})
	if err != nil {
		return &types.FriendListResp{}, err
	}
	userRecords := make(map[string]*userclient.UserEntity, len(user.Users))
	for i, _ := range user.Users {
		userRecords[user.Users[i].Id] = user.Users[i]
	}

	// 组装数据并返回
	respList := make([]*types.Friends, 0, len(friends.List))
	for _, f := range friends.List {
		friend := &types.Friends{
			Id:        f.Id,
			FriendUid: f.FriendUid,
		}
		if userRecord, ok := userRecords[f.FriendUid]; ok {
			friend.Avatar = userRecord.Avatar
			friend.Nickname = userRecord.Nickname
		}
		respList = append(respList, friend)
	}

	return &types.FriendListResp{
		List: respList,
	}, nil
}
