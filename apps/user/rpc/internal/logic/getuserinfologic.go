package logic

import (
	"app/apps/user/models"
	"app/apps/user/rpc/internal/svc"
	"app/apps/user/rpc/user"
	"github.com/pkg/errors"

	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = errors.New("用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserInfo 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {

	userEntity, err := l.svcCtx.UserModels.FindOne(l.ctx, in.Id)
	if err != nil {
		if nil == models.ErrNotFound {
			return nil, errors.Wrapf(ErrUserNotFound, "用户不存在: %v", in.Id)
		} else {
			return nil, errors.WithStack(err)
		}
	}

	var resp user.UserEntity
	err = copier.Copy(&resp, userEntity)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &user.GetUserInfoResponse{
		User: &resp,
	}, nil
}
