package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupMembersModel = (*customGroupMembersModel)(nil)

type (
	// GroupMembersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupMembersModel.
	GroupMembersModel interface {
		groupMembersModel
		TransCtxInsert(ctx context.Context, session sqlx.Session, data *GroupMembers) (sql.Result, error)
	}

	customGroupMembersModel struct {
		*defaultGroupMembersModel
	}
)

// TransCtxInsert 支持事务的insert
func (c *customGroupMembersModel) TransCtxInsert(ctx context.Context, session sqlx.Session, data *GroupMembers) (sql.Result, error) {
	groupMembersIdKey := fmt.Sprintf("%s%v", cacheGroupMembersIdPrefix, data.Id)
	ret, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", c.table, groupMembersRowsExpectAutoSet)
		// 如果有 session 的话，就用这个 session 来执行操作
		if session != nil {
			return session.ExecCtx(ctx, query, data.GroupId, data.UserId, data.RoleLevel, data.JoinTime, data.JoinSource, data.InviterUid, data.OperatorUid)
		}
		return conn.ExecCtx(ctx, query, data.GroupId, data.UserId, data.RoleLevel, data.JoinTime, data.JoinSource, data.InviterUid, data.OperatorUid)
	}, groupMembersIdKey)
	return ret, err
}

// NewGroupMembersModel returns a model for the database table.
func NewGroupMembersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupMembersModel {
	return &customGroupMembersModel{
		defaultGroupMembersModel: newGroupMembersModel(conn, c, opts...),
	}
}
