package logic

import (
	socialmodels "app/apps/social/models"
	"app/pkg/xerr"
	"context"
	"github.com/pkg/errors"
	"time"

	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutIn 提交加群申请
func (l *GroupPutInLogic) GroupPutIn(in *rpc.GroupPutInReq) (*rpc.GroupPutInResp, error) {
	if in.ReqTime == 0 {
		in.ReqTime = time.Now().Unix()
	}

	//  1. 普通用户申请 ： 如果群无验证直接进入
	//  2. 群成员邀请： 如果群无验证直接进入
	//  3. 群管理员/群创建者邀请：直接进入群

	var (
		groupMembers *socialmodels.GroupMembers
	)
	// 1.判断用户是否已经在群中
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByGroupIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBError(), "find group member by groud id and req id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if groupMembers != nil {
		return nil, errors.Wrapf(xerr.NewMsgErr("用户已加入该群"), "reqId: %s, groupId: %s", in.ReqId, in.GroupId)
	}

	// 2.判断用户的加群申请是否已经存在 todo
	//l.svcCtx.GroupRequestsModel.FindOne()

	return &rpc.GroupPutInResp{}, nil
}
