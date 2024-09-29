package logic

import (
	socialmodels "app/apps/social/models"
	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"
	"app/pkg/constants"
	"app/pkg/xerr"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutIn 添加好友申请
func (l *FriendPutInLogic) FriendPutIn(in *rpc.FriendPutInReq) (*rpc.FriendPutInResp, error) {
	// 请求时间不为 0 时取当前时间
	if in.ReqTime == 0 {
		in.ReqTime = time.Now().Unix()
	}
	// 1.判断申请人与目标好友是否已经是好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFriendUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBError(), "find frends err by userId: %s, reqUid: %s", in.UserId, in.ReqUid)
	}
	if friends != nil {
		// 好友关系已经存在
		return nil, errors.Wrapf(xerr.NewMsgErr("好友关系已存在"), "userId: %s, reqUid: %s", in.UserId, in.ReqUid)
	}

	// 2.判断是否已经存在好友申请记录
	friendRequest, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBError(), "find friend request err by reqUid: %s, userId: %s", in.ReqUid, in.UserId)
	}
	if friendRequest != nil {
		// 好友申请记录已存在
		return nil, errors.Wrapf(xerr.NewMsgErr("好友申请记录已存在"), "reqUid: %s, userId: %s", in.ReqUid, in.UserId)
	}

	// 3.添加好友申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: in.ReqTime, // 添加好友的时间只精确到秒
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult), // 添加好友请求创建时默认是未处理
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "insert friend request err by reqUid: %s, userId: %s", in.ReqUid, in.UserId)
	}

	return &rpc.FriendPutInResp{}, nil
}
