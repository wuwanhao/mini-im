package group

import (
	"app/apps/social/rpc/socialclient"
	"app/apps/user/rpc/userclient"
	"context"

	"app/apps/social/api/internal/svc"
	"app/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 获取某个群的成员列表
func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {

	// 获取群成员列表的所有用户 ID
	memberList, err := l.svcCtx.SocialRpc.GroupMemberList(l.ctx, &socialclient.GroupMemberListReq{GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	uidList := make([]string, 0, len(memberList.List))
	for _, member := range memberList.List {
		uidList = append(uidList, member.UserId)
	}

	// 获取所有用户的详细信息
	userList, err := l.svcCtx.UserRpc.FindUser(l.ctx, &userclient.FindUserRequest{
		Ids: uidList,
	})
	userInfoList := make(map[string]*userclient.UserEntity, len(userList.Users))
	for i, _ := range userList.Users {
		userInfoList[userList.Users[i].Id] = userList.Users[i]
	}

	// 组装数据并返回
	respList := make([]*types.GroupMembers, 0, len(memberList.List))
	for _, u := range memberList.List {
		member := &types.GroupMembers{
			GroupId:    u.GroupId,
			Id:         int64(u.Id),
			InviterUid: u.InviterUid,
			RoleLevel:  int(u.RoleLevel),
		}

		if v, ok := userInfoList[u.UserId]; ok {
			member.Nickname = v.Nickname
			member.UserAvatarUrl = v.Avatar
		}
		respList = append(respList, member)
	}

	return &types.GroupUserListResp{
		List: respList,
	}, nil

}
