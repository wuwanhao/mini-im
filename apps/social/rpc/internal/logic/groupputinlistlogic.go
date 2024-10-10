package logic

import (
	"app/pkg/xerr"
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取某个群的加群申请列表
func (l *GroupPutInListLogic) GroupPutInList(in *rpc.GroupPutInListReq) (*rpc.GroupPutInListResp, error) {

	groupRequests, err := l.svcCtx.GroupRequestsModel.ListNoHanderByGroupId(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "list group req err:%v, req:%+v", err, in)
	}
	var groupRequestsList []*rpc.GroupRequests
	copier.Copy(&groupRequestsList, groupRequests)

	return &rpc.GroupPutInListResp{
		List: groupRequestsList,
	}, nil
}
