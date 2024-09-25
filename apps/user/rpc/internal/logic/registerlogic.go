package logic

import (
	"app/apps/user/models"
	"app/apps/user/rpc/internal/svc"
	"app/apps/user/rpc/user"
	"app/pkg/ctxdata"
	"app/pkg/encrypt"
	"app/pkg/wuid"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegistered = errors.New("手机号已注册")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 1.检查手机号是否已经注册过
	userEntity, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}
	if userEntity != nil {
		return nil, ErrPhoneIsRegistered
	}

	uid := wuid.GenUID(l.svcCtx.Config.Mysql.Datasource)
	fmt.Println("uid:", uid)
	// 2.定义用户数据
	userEntity = &models.Users{
		Id:       uid,
		Phone:    in.Phone,
		Nickname: in.Nickname,
		Avatar:   in.Avatar,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	// 3.处理用户的密码
	if len(in.Password) > 0 {
		hashedPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		}
	}

	// 4.保存到数据库
	_, err = l.svcCtx.UserModels.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// 5.生成返回token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}
	return &user.RegisterResponse{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
