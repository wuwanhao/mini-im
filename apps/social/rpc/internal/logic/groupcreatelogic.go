package logic

import (
	"context"

	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群组业务：创建群组、修改群、群公告、申请群、加群请求列表、加群请求处理...
func (l *GroupCreateLogic) GroupCreate(in *rpc.GroupCreateReq) (*rpc.GroupCreateResp, error) {
	// todo: add your logic here and delete this line

	return &rpc.GroupCreateResp{}, nil
}
