package logic

import (
	socialmodels "app/apps/social/models"
	"app/pkg/constants"
	"app/pkg/xerr"
	"context"
	"database/sql"
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

	//  1. 普通用户申请：如果群无验证直接进入
	//  2. 群成员邀请：如果群无验证直接进入
	//  3. 群管理员/群创建者邀请：直接进入群

	var (
		groupMembers *socialmodels.GroupMembers
		groupInfo    *socialmodels.Groups
	)
	// 1.判断用户是否已经在群中
	groupMembers, err := l.svcCtx.GroupMembersModel.FindByGroupIdAndUserId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBError(), "find group member by groud id and req id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if groupMembers != nil {
		return nil, errors.Wrapf(xerr.NewMsgErr("用户已加入该群"), "reqId: %s, groupId: %s", in.ReqId, in.GroupId)
	}

	// 2.判断用户的加群申请是否已经存在
	groupRequests, err := l.svcCtx.GroupRequestsModel.FindByReqIdAndGroupId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBError(), "find group request by req id and group id err %v, req %v, %v", err,
			in.ReqId, in.GroupId)
	}
	if groupRequests != nil {
		return nil, errors.Wrapf(xerr.NewMsgErr("已存在该用户的加群申请"), "reqId: %s, groupId: %s", in.ReqId, in.GroupId)
	}

	// 3.若加群申请不存在，则构建一条加群申请
	groupRequests = &socialmodels.GroupRequests{
		GroupId: in.GroupId,
		ReqId:   in.ReqId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: in.ReqTime,
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}

	// 创建群用户的匿名函数
	createGroupMember := func() {
		if err != nil {
			return
		}
		err = l.CreateGroupMember(in)
	}

	// 查询要申请加入的群是否存在
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "find group by group id err %v, req %v", err, in)
	}

	// 检查加群是否需要验证
	if !groupInfo.IsVerify {
		// 不需要验证
		defer createGroupMember()
		groupRequests.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupRequests.HandleTime = time.Now().Unix()
		return l.CreateGroupReq(groupRequests, true)
	}

	// 进群方式验证: 申请，直接创建一个请求，然后返回响应，等待审核处理
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		return l.CreateGroupReq(groupRequests, false)
	}

	// 进群方式验证: 邀请，记录邀请人信息，直接创建群成员，然后返回响应
	// 获取当前邀请人信息
	inviter, err := l.svcCtx.GroupMembersModel.FindByGroupIdAndInvitorId(l.ctx, in.GroupId, in.InviterUid)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "find group member by group id and invitor id err %v, req %v", err, in)
	}

	// 如果邀请人信息是群主或者管理员
	if (constants.GroupRoleLevel(inviter.RoleLevel) == constants.CreatorGroupRoleLevel) || (constants.
		GroupRoleLevel(inviter.RoleLevel) == constants.ManagerGroupRoleLevel) {
		// 直接创建群成员
		defer createGroupMember()
		groupRequests.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupRequests.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		groupRequests.HandleTime = time.Now().Unix()
		return l.CreateGroupReq(groupRequests, true)
	}

	return l.CreateGroupReq(groupRequests, true)
}

// 创建一个加群请求
func (l *GroupPutInLogic) CreateGroupReq(groupReq *socialmodels.GroupRequests, isPass bool) (*rpc.GroupPutInResp, error) {
	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "insert group request err %v, req %v", err, groupReq)
	}

	if isPass {
		return &rpc.GroupPutInResp{GroupId: groupReq.GroupId}, nil
	}
	return &rpc.GroupPutInResp{}, nil

}

// 创建群成员
func (l *GroupPutInLogic) CreateGroupMember(in *rpc.GroupPutInReq) error {
	groupMembers := &socialmodels.GroupMembers{
		GroupId:   in.GroupId,
		UserId:    in.ReqId,
		RoleLevel: int64(constants.AtLargeGroupRoleLevel),
		OperatorUid: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		JoinTime: time.Now().Unix(),
	}
	_, err := l.svcCtx.GroupMembersModel.TransCtxInsert(l.ctx, nil, groupMembers)
	if err != nil {
		return errors.Wrapf(xerr.NewDBError(), "insert group member err %v, req %v", err, in)
	}
	return nil
}
