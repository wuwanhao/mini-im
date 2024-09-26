package user

import (
	"app/apps/user/rpc/userclient"
	"context"

	"app/apps/user/api/internal/svc"
	"app/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 用户登录
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	// l.svcCtx.User 即为 grpc 客户端
	loginResponse, err := l.svcCtx.User.Login(l.ctx, &userclient.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// 构造返回内容
	var res types.LoginResp
	res.Token = loginResponse.Token
	res.Expire = loginResponse.Expire
	return &res, nil
}
