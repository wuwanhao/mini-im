package logic

import (
	"app/pkg/constants"
	"context"
	"github.com/pkg/errors"

	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"

	"github.com/zeromicro/go-zero/core/logx"
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
		return nil, errors.Wrapf(err, "get UnHandled FriendRequest error, ReqId:%d", in.FriendReqId)
	}
	// 2.判断当前记录的处理状态
	switch constants.HandlerResult(currentFriendRequests.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.Wrapf(errors.New("already handled"), "already handled")
	case constants.RefuseHandlerResult:
		return nil, errors.Wrapf(errors.New("already handled"), "already handled")
	}

	// 3.处理（事务）
	currentFriendRequests.HandleResult.Int64 = int64(in.HandleResult)
	// 3.1 更新好友申请记录状态 todo
	l.svcCtx.FriendRequestsModel.Trans()
	// 3.2 添加好友关系

	return &rpc.FriendPutInHandleResp{}, nil
}
