package svc

import (
	socialmodels "app/apps/social/models"
	"app/apps/social/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 将用到的模型在服务核心上下文中进行加载
	socialmodels.GroupsModel
	socialmodels.GroupMembersModel
	socialmodels.GroupRequestsModel
	socialmodels.FriendsModel
	socialmodels.FriendRequestsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:              c,
		GroupsModel:         socialmodels.NewGroupsModel(sqlConn, c.Cache),
		GroupMembersModel:   socialmodels.NewGroupMembersModel(sqlConn, c.Cache),
		GroupRequestsModel:  socialmodels.NewGroupRequestsModel(sqlConn, c.Cache),
		FriendsModel:        socialmodels.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: socialmodels.NewFriendRequestsModel(sqlConn, c.Cache),
	}
}
