package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ FriendsModel = (*customFriendsModel)(nil)

type (
	// FriendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendsModel.
	FriendsModel interface {
		friendsModel
		FindByUidAndFriendUid(ctx context.Context, uid string, friendUid string) (*Friends, error)
		TransBatchInsert(ctx context.Context, session sqlx.Session, data ...*Friends) (sql.Result, error)
	}

	customFriendsModel struct {
		*defaultFriendsModel
	}
)

// FindByUidAndFriendUid 通过 user_id 和 friend_uid 查询是否是好友关系
func (c *customFriendsModel) FindByUidAndFriendUid(ctx context.Context, uid string, friendUid string) (*Friends, error) {

	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `friend_uid` = ?", friendsRows, c.table)
	var friends Friends
	err := c.QueryRowNoCacheCtx(ctx, &friends, query, uid, friendUid)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &friends, nil
}

// TransBatchInsert 支持事务的批量插入
func (m *defaultFriendsModel) TransBatchInsert(ctx context.Context, session sqlx.Session, data ...*Friends) (sql.Result, error) {
	// 生成 mysql 的 批量插入语句
	var (
		builder strings.Builder
		args    []any
	)
	if len(data) == 0 {
		return nil, nil
	}
	builder.WriteString(fmt.Sprintf("insert into %s (%s) values ", m.table, friendsRowsExpectAutoSet))
	for i, v := range data {
		builder.WriteString("(?, ?, ?, ?)")
		args = append(args, v.UserId, v.FriendUid, v.Remark, v.AddSource)
		if i != len(data)-1 {
			builder.WriteString(",")
		}
	}
	fmt.Println(builder.String())
	return session.ExecCtx(ctx, builder.String(), args...)
}

// NewFriendsModel returns a model for the database table.
func NewFriendsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendsModel {
	return &customFriendsModel{
		defaultFriendsModel: newFriendsModel(conn, c, opts...),
	}
}
