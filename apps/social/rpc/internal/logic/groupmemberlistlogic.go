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

type GroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberListLogic {
	return &GroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupMemberList 获取群里的所有群成员
func (l *GroupMemberListLogic) GroupMemberList(in *rpc.GroupMemberListReq) (*rpc.GroupMemberListResp, error) {
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByGroupId(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "list group members error, groupId:%+v, err:%+v", in.GroupId, err)
	}
	if len(groupMembers) == 0 {
		return &rpc.GroupMemberListResp{}, nil
	}

	var resp []*rpc.GroupMembers
	copier.Copy(&resp, groupMembers)

	return &rpc.GroupMemberListResp{
		List: resp,
	}, nil
}
