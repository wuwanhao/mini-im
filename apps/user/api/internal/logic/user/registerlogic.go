package user

import (
	"app/apps/user/rpc/userclient"
	"context"

	"app/apps/user/api/internal/svc"
	"app/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 注册
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerResponse, err := l.svcCtx.User.Register(l.ctx, &userclient.RegisterRequest{
		Phone:    req.Phone,
		Password: req.Password,
		Nickname: req.NickName,
		Sex:      int32(req.Sex),
		Avatar:   req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		Token:  registerResponse.Token,
		Expire: registerResponse.Expire,
	}, nil

}
