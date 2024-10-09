package models

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupRequestsModel = (*customGroupRequestsModel)(nil)

type (
	// GroupRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupRequestsModel.
	GroupRequestsModel interface {
		groupRequestsModel
		FindByReqIdAndGroupId(ctx context.Context, reqId, groupId string) (*GroupRequests, error)
	}

	customGroupRequestsModel struct {
		*defaultGroupRequestsModel
	}
)

// FindByReqIdAndGroupId todo
func (c *defaultGroupRequestsModel) FindByReqIdAndGroupId(ctx context.Context, reqId, groupId string) (*GroupRequests, error) {
	return nil, nil
}

// NewGroupRequestsModel returns a model for the database table.
func NewGroupRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestsModel {
	return &customGroupRequestsModel{
		defaultGroupRequestsModel: newGroupRequestsModel(conn, c, opts...),
	}
}
