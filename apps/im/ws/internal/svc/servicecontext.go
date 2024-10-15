package svc

import (
	"app/apps/im/ws/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(conf config.Config) *ServiceContext {
	return &ServiceContext{
		Config: conf,
	}
}
