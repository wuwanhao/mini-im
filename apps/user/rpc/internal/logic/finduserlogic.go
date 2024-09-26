package logic

import (
	"app/apps/user/models"
	"app/apps/user/rpc/internal/svc"
	"app/apps/user/rpc/user"
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserRequest) (*user.FindUserResponse, error) {
	var (
		userEntities []*models.Users
		err          error
	)

	// 条件查询
	if in.Phone != "" {
		users, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntities = append(userEntities, users)
		}
	} else if in.Name != "" {
		userEntities, err = l.svcCtx.UserModels.FindUserByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntities, err = l.svcCtx.UserModels.FindUserByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	var resp []*user.UserEntity
	err = copier.Copy(&resp, userEntities)
	if err != nil {
		return nil, err
	}
	return &user.FindUserResponse{
		Users: resp,
	}, nil
}
