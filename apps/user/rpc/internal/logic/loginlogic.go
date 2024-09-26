package logic

import (
	"app/apps/user/models"
	"app/apps/user/rpc/internal/svc"
	"app/apps/user/rpc/user"
	"app/pkg/ctxdata"
	"app/pkg/encrypt"
	"app/pkg/xerr"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var (

	// 使用系统内部自定义的错误对象
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号未注册")
	ErrPasswordError    = xerr.New(xerr.SERVER_COMMON_ERROR, "密码错误")
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
	if err != nil {
		if err == models.ErrNotFound {
			// errors.WithStack(ErrPhoneNotRegister) 返回普通形式的错误
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		//  errors.Wrapf(xerr.NewDBError(), "findByPhone err:%v, param:%v", err, in.Phone) 返回详细信息的错误
		return nil, errors.Wrapf(xerr.NewDBError(), "findByPhone err:%v, param:%v", err, in.Phone)
	}

	// 2.校验密码
	validateResult := encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String)
	if !validateResult {
		return nil, errors.WithStack(ErrPasswordError)
	}

	// 3.生成并返回 token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBError(), "GetJwtToken err:%v", err)
	}
	return &user.LoginResponse{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil

}
