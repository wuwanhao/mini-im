package group

import (
	"app/apps/social/rpc/socialclient"
	"context"
	"github.com/jinzhu/copier"

	"app/apps/social/api/internal/svc"
	"app/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 获取某个群的加群列表
func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListRep) (resp *types.GroupPutInListResp, err error) {
	list, err := l.svcCtx.SocialRpc.GroupPutInList(l.ctx, &socialclient.GroupPutInListReq{
		GroupId: req.GroupId,
	})

	var respList []*types.GroupRequests
	copier.Copy(&respList, list.List)

	return &types.GroupPutInListResp{List: respList}, nil

	return
}
