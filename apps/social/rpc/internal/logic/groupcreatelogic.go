package logic

import (
	socialmodels "app/apps/social/models"
	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"
	"app/pkg/constants"
	"app/pkg/wuid"
	"app/pkg/xerr"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

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

	// 2.执行添加操作，事务，创建群的同时创建群成员
	groups := &socialmodels.Groups{
		Id:         wuid.GenUID(l.svcCtx.Config.Mysql.DataSource),
		Name:       in.Name,
		Icon:       in.Icon,
		CreatorUid: in.CreatorUid,
		Status: sql.NullInt64{
			Int64: int64(in.Status),
			Valid: true,
		},
		IsVerify: false,
	}
	// 2.1 事务开始
	err = l.svcCtx.GroupsModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 插入 Groups
		_, err = l.svcCtx.GroupsModel.TransCtxInsert(l.ctx, session, groups)
		if err != nil {
			return errors.Wrapf(xerr.NewDBError(), "insert group err by name: %s, err: %v", in.Name, err)
		}

		// 插入 GroupsMembers
		_, err = l.svcCtx.GroupMembersModel.TransCtxInsert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,
			UserId:    groups.CreatorUid,
			RoleLevel: int64(constants.CreatorGroupRoleLevel),
			JoinTime:  time.Now().Unix(),
			JoinSource: sql.NullInt64{
				Int64: int64(constants.PutInGroupJoinSource),
				Valid: true,
			},
			OperatorUid: sql.NullString{
				String: in.CreatorUid,
				Valid:  true,
			},
		})

		if err != nil {
			return errors.Wrapf(xerr.NewDBError(), "insert group users err by name: %s, err: %v", in.Name, err)
		}

		return nil
	})
	return &rpc.GroupCreateResp{}, err
}
