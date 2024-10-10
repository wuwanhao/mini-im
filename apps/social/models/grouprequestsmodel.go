package models

import (
	"context"
	"fmt"
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

// FindByReqIdAndGroupId 根据 reqID 和 groupId 查询 用户的加群申请
func (c *defaultGroupRequestsModel) FindByReqIdAndGroupId(ctx context.Context, reqId, groupId string) (*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_id` = ? and `group_id` = ?", groupRequestsRows, c.table)
	var resp *GroupRequests
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, reqId, groupId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NewGroupRequestsModel returns a model for the database table.
func NewGroupRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestsModel {
	return &customGroupRequestsModel{
		defaultGroupRequestsModel: newGroupRequestsModel(conn, c, opts...),
	}
}
