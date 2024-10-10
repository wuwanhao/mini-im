package logic

import (
	socialmodels "app/apps/social/models"
	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"
	"app/pkg/constants"
	"app/pkg/ctxdata"
	"app/pkg/xerr"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrGroupReqBeforePass   = xerr.NewMsgErr("请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsgErr("请求已拒绝")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 处理加群申请
func (l *GroupPutInHandleLogic) GroupPutInHandle(in *rpc.GroupPutInHandleReq) (*rpc.GroupPutInHandleResp, error) {
	// 找到这条加群申请
	groupRequests, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, int64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "find group request by id err %v, id %v", err, in.GroupReqId)
	}

	// 当前加群申请的状态，如果处于终态，则报错
	switch constants.HandlerResult(groupRequests.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}

	groupRequests.HandleTime = time.Now().Unix()
	groupRequests.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}
	groupRequests.HandleUserId = sql.NullString{
		String: ctxdata.GetUid(l.ctx),
		Valid:  true,
	}

	// 事务开启
	l.svcCtx.GroupRequestsModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1.更新加群申请状态
		if err := l.svcCtx.GroupRequestsModel.TransCtxUpdate(l.ctx, session, groupRequests); err != nil {
			return errors.Wrapf(xerr.NewDBError(), "update group request err %v, groupReqId %v", err, in.GroupReqId)
		}
		// 2.若申请通过，则添加群-成员关系，否则不添加，直接返回响应
		if constants.HandlerResult(in.HandleResult) == constants.RefuseHandlerResult {
			return nil
		}
		groupMembers := &socialmodels.GroupMembers{
			GroupId:     groupRequests.GroupId,
			UserId:      groupRequests.ReqId,
			JoinTime:    time.Now().Unix(),
			RoleLevel:   int64(constants.AtLargeGroupRoleLevel),
			JoinSource:  groupRequests.JoinSource,
			InviterUid:  groupRequests.InviterUserId,
			OperatorUid: groupRequests.HandleUserId,
		}

		_, err = l.svcCtx.GroupMembersModel.TransCtxInsert(l.ctx, session, groupMembers)
		if err != nil {
			return errors.Wrapf(xerr.NewDBError(), "insert group members err %v, groupMembers %v", err, groupMembers)
		}

		return nil
	})

	return &rpc.GroupPutInHandleResp{}, nil
}
