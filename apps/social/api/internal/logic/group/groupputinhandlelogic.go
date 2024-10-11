package group

import (
	"app/apps/social/rpc/socialclient"
	"app/pkg/constants"
	"app/pkg/ctxdata"
	"context"

	"app/apps/social/api/internal/svc"
	"app/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutInHandle 加群请求处理
func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	uid := ctxdata.GetUid(l.ctx)

	_, err = l.svcCtx.SocialRpc.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,
		HandleResult: req.HandleResult,
		HandleUid:    uid,
		GroupId:      req.GroupId,
	})
	if err != nil {
		return nil, err
	}

	if constants.HandlerResult(req.HandleResult) != constants.PassHandlerResult {
		return
	}

	return
}
