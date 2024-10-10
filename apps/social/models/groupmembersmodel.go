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
		ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error)
		ListByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error)
		FindByGroupIdAndUserId(ctx context.Context, groupId, userId string) (*GroupMembers, error)
		FindByGroupIdAndInvitorId(ctx context.Context, groupId, invitorId string) (*GroupMembers, error)
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

// ListByUserId 根据 userId 获取 群-用户 列表
func (c *customGroupMembersModel) ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", groupMembersRows, c.table)
	var resp []*GroupMembers
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListByGroupId 根据 groupId 获取 群-用户 列表
func (c *customGroupMembersModel) ListByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ?", groupMembersRows, c.table)
	var resp []*GroupMembers
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindByGroupIdAndUserId 根据 groupId 和 userId 获取 群-用户 列表
func (c *customGroupMembersModel) FindByGroupIdAndUserId(ctx context.Context, groupId, userId string) (*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `user_id` = ?", groupMembersRows, c.table)
	var resp *GroupMembers
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, groupId, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindByGroupIdAndUserId 根据 groupId 和 invitorId 获取 群-用户 列表
func (c *customGroupMembersModel) FindByGroupIdAndInvitorId(ctx context.Context, groupId, InvitorId string) (*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `inviter_uid` = ?", groupMembersRows, c.table)
	var resp *GroupMembers
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, groupId, InvitorId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NewGroupMembersModel returns a model for the database table.
func NewGroupMembersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupMembersModel {
	return &customGroupMembersModel{
		defaultGroupMembersModel: newGroupMembersModel(conn, c, opts...),
	}
}
