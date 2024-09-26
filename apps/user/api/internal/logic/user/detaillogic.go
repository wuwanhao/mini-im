package user

import (
	"app/apps/user/rpc/userclient"
	"app/pkg/ctxdata"
	"context"
	"github.com/jinzhu/copier"

	"app/apps/user/api/internal/svc"
	"app/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context     // 应用程序上下文
	svcCtx *svc.ServiceContext // 核心服务上下文
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Detail 获取用户详情
func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// 先在上下文中获取用户 ID
	uid := ctxdata.GetUid(l.ctx)

	userInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &userclient.GetUserInfoRequest{Id: uid})
	if err != nil {
		return nil, err
	}

	var user types.User
	err = copier.Copy(&user, userInfoResp.User)
	if err != nil {
		return nil, err
	}
	return &types.UserInfoResp{Info: user}, nil

}
