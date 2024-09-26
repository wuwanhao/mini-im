package svc

import (
	"app/apps/user/api/internal/config"
	"app/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	userclient.User // 引用 grpc 客户端
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
