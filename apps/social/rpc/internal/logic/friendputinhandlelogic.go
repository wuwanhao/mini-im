package logic

import (
	socialmodels "app/apps/social/models"
	"app/pkg/constants"
	"app/pkg/xerr"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsgErr("好友申请已经通过")
	ErrFriendReqBeforeRefuse = xerr.NewMsgErr("好友申请已经拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutInHandle 处理某条好友申请
func (l *FriendPutInHandleLogic) FriendPutInHandle(in *rpc.FriendPutInHandleReq) (*rpc.FriendPutInHandleResp, error) {
	// 1.获取当前用户待处理的好友申请记录
	currentFriendRequests, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, int64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "get UnHandled FriendRequest error, ReqId:%d", in.FriendReqId)
	}
	// 2.判断当前记录的处理状态
	switch constants.HandlerResult(currentFriendRequests.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}

	// 3.处理（事务开始）
	currentFriendRequests.HandleResult.Int64 = int64(in.HandleResult)
	fmt.Println(currentFriendRequests.HandleResult.Int64)
	fmt.Println("-----------")
	err = l.svcCtx.FriendRequestsModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		//  3.1 修改好友关系状态
		if err := l.svcCtx.FriendRequestsModel.TransUpdate(ctx, session, currentFriendRequests); err != nil {
			return errors.Wrapf(xerr.NewDBError(), "Update FriendRequests error, ReqId:%v, err:%v", in.FriendReqId, err)
		}

		// 3.2 添加好友关系，建立两条好友关系记录，因为friends 表是冗余设计，所以添加两条记录
		if in.HandleResult == int32(constants.PassHandlerResult) {
			friends := []*socialmodels.Friends{
				{
					UserId:    currentFriendRequests.ReqUid,
					FriendUid: currentFriendRequests.UserId,
				},
				{
					UserId:    currentFriendRequests.UserId,
					FriendUid: currentFriendRequests.ReqUid,
				},
			}
			_, err := l.svcCtx.FriendsModel.TransBatchInsert(ctx, session, friends...)
			if err != nil {
				return errors.Wrapf(xerr.NewDBError(), "BatchInsert Friends error, ReqId:%v, err:%v", in.FriendReqId, err)
			}
		}
		return nil
	})

	return &rpc.FriendPutInHandleResp{}, err
}
