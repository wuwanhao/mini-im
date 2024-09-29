package logic

import (
	"context"

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

// todo 提交加群申请
func (l *GroupPutInLogic) GroupPutIn(in *rpc.GroupPutInReq) (*rpc.GroupPutInResp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("GroupPutIn", in)

	return &rpc.GroupPutInResp{}, nil
}
