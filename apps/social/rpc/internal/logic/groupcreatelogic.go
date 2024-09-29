package logic

import (
	socialmodels "app/apps/social/models"
	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"
	"app/pkg/wuid"
	"app/pkg/xerr"
	"context"
	"database/sql"
	"github.com/pkg/errors"

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

// GroupCreate 创建群组
func (l *GroupCreateLogic) GroupCreate(in *rpc.GroupCreateReq) (*rpc.GroupCreateResp, error) {

	l.Logger.Info("GroupCreate", in)
	// 1.检查是否有同名群
	group, err := l.svcCtx.GroupsModel.FindOneByName(l.ctx, in.Name)
	if group != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBError(), "group name repeat by name: %s", in.Name)
	}

	// 2.todo: 这里一些步骤需要进一步完善 执行添加操作
	_, err = l.svcCtx.GroupsModel.Insert(l.ctx, &socialmodels.Groups{
		Id:   wuid.GenUID(l.svcCtx.Config.Mysql.DataSource),
		Name: in.Name,
		Icon: in.Icon,
		Status: sql.NullInt64{
			Int64: int64(in.Status),
			Valid: true,
		},
		CreatorUid: in.CreatorUid,
		GroupType:  0,
		IsVerify:   false,
		Notification: sql.NullString{
			String: "",
			Valid:  true,
		},
		NotificationUid: sql.NullString{
			String: "",
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "insert group err by name: %s, err: %v", in.Name, err)
	}

	return &rpc.GroupCreateResp{}, nil
}
