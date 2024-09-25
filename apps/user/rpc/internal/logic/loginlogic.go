package logic

import (
	"app/apps/user/models"
	"app/apps/user/rpc/internal/svc"
	"app/apps/user/rpc/user"
	"app/pkg/ctxdata"
	"app/pkg/encrypt"
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var (
	ErrPhoneNotRegister = errors.New("手机号未注册")
	ErrPasswordError    = errors.New("密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 用户登录
func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {

	// 1.通过手机号验证用户是否注册
	userEntity, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
	fmt.Println("userEntity:", userEntity)
	fmt.Println(in.Phone)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, ErrPhoneNotRegister
		}
		return nil, err
	}

	// 2.校验密码
	validateResult := encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String)
	if !validateResult {
		return nil, ErrPasswordError
	}

	// 3.生成并返回 token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}
	return &user.LoginResponse{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil

}
