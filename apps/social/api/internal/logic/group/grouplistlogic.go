package group

import (
	"app/apps/social/rpc/socialclient"
	"app/pkg/ctxdata"
	"context"
	"github.com/jinzhu/copier"

	"app/apps/social/api/internal/svc"
	"app/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 获取某个用户的	群列表
func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	groupList, err := l.svcCtx.SocialRpc.GroupList(l.ctx, &socialclient.GroupListReq{UserId: uid})
	if err != nil {
		return nil, err
	}

	var list []*types.Groups
	copier.Copy(&list, groupList.List)
	return &types.GroupListResp{List: list}, nil
}
