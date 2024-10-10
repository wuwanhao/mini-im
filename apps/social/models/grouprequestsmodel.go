package models

import (
	"app/pkg/constants"
	"context"
	"database/sql"
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
		ListNoHanderByGroupId(ctx context.Context, groupId string) ([]*GroupRequests, error)
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		TransCtxUpdate(ctx context.Context, session sqlx.Session, data *GroupRequests) error
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

func (c *defaultGroupRequestsModel) ListNoHanderByGroupId(ctx context.Context, groupId string) ([]*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `handle_result` = ?", groupRequestsRows, c.table)
	var resp []*GroupRequests
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, groupId, constants.NoHandlerResult)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Trans 用于执行事务的方法
func (c *defaultGroupRequestsModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return c.TransactCtx(ctx, func(ctx context.Context, conn sqlx.Session) error {
		return fn(ctx, conn)
	})
}

// 支持事务的 更新操作
func (m *defaultGroupRequestsModel) TransCtxUpdate(ctx context.Context, session sqlx.Session, data *GroupRequests) error {
	groupRequestsIdKey := fmt.Sprintf("%s%v", cacheGroupRequestsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupRequestsRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.ReqId, data.GroupId, data.ReqMsg, data.ReqTime, data.JoinSource, data.InviterUserId, data.HandleUserId, data.HandleTime, data.HandleResult, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.ReqId, data.GroupId, data.ReqMsg, data.ReqTime, data.JoinSource, data.InviterUserId, data.HandleUserId, data.HandleTime, data.HandleResult, data.Id)
	}, groupRequestsIdKey)
	return err
}

// NewGroupRequestsModel returns a model for the database table.
func NewGroupRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestsModel {
	return &customGroupRequestsModel{
		defaultGroupRequestsModel: newGroupRequestsModel(conn, c, opts...),
	}
}
