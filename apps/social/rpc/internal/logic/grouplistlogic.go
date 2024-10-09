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

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupList 获取某个用户的所有群列表
func (l *GroupListLogic) GroupList(in *rpc.GroupListReq) (*rpc.GroupListResp, error) {
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "list group memebers error,  userId:%s err:%+v", in.UserId, err)
	}
	if len(groupMembers) == 0 {
		return &rpc.GroupListResp{}, nil
	}

	// 拿到所有的 groupId
	ids := make([]string, 0, len(groupMembers))
	for _, member := range groupMembers {
		ids = append(ids, member.GroupId)
	}

	// 根据 groupId 拿到所有的群信息
	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "list groups error, groupIds:%+v, err:%+v", ids, err)
	}
	var resp []*rpc.Groups
	copier.Copy(&resp, &groups)
	return &rpc.GroupListResp{
		List: resp,
	}, nil
}
