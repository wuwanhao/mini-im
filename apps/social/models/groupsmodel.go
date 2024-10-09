package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupsModel = (*customGroupsModel)(nil)

type (
	// GroupsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupsModel.
	GroupsModel interface {
		groupsModel
		FindOneByName(ctx context.Context, name string) (*Groups, error)
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		TransCtxInsert(ctx context.Context, session sqlx.Session, data *Groups) (sql.Result, error)
	}

	customGroupsModel struct {
		*defaultGroupsModel
	}
)

// FindOneByName 根据群名查找群
func (c *customGroupsModel) FindOneByName(ctx context.Context, name string) (*Groups, error) {
	query := fmt.Sprintf("select %s from %s where `name` = ?", groupsRows, c.table)
	var groups Groups
	err := c.QueryRowNoCacheCtx(ctx, &groups, query, name)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}

	return &groups, nil
}

// Trans 用于执行事务的方法
func (c *customGroupsModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return c.TransactCtx(ctx, func(ctx context.Context, conn sqlx.Session) error {
		return fn(ctx, conn)
	})
}

// TransCtxInsert 支持事务的插入
func (m *defaultGroupsModel) TransCtxInsert(ctx context.Context, session sqlx.Session, data *Groups) (sql.Result, error) {
	groupsIdKey := fmt.Sprintf("%s%v", cacheGroupsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupsRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.Id, data.Name, data.Icon, data.Status, data.CreatorUid, data.GroupType, data.IsVerify, data.Notification, data.NotificationUid)
		}
		return conn.ExecCtx(ctx, query, data.Id, data.Name, data.Icon, data.Status, data.CreatorUid, data.GroupType, data.IsVerify, data.Notification, data.NotificationUid)
	}, groupsIdKey)
	return ret, err
}

// NewGroupsModel returns a model for the database table.
func NewGroupsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupsModel {
	return &customGroupsModel{
		defaultGroupsModel: newGroupsModel(conn, c, opts...),
	}
}
