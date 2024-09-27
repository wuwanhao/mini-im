package models

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendsModel = (*customFriendsModel)(nil)

type (
	// FriendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendsModel.
	FriendsModel interface {
		friendsModel
		FindByUidAndFriendUid(ctx context.Context, uid string, friendUid string) (*Friends, error)
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

// NewFriendsModel returns a model for the database table.
func NewFriendsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendsModel {
	return &customFriendsModel{
		defaultFriendsModel: newFriendsModel(conn, c, opts...),
	}
}
