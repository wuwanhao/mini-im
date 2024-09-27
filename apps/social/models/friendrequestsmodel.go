package models

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendRequestsModel = (*customFriendRequestsModel)(nil)

type (
	// FriendRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestsModel.
	FriendRequestsModel interface {
		friendRequestsModel
		FindByReqUidAndUserId(ctx context.Context, reqUid string, userId string) (*FriendRequests, error)
	}

	customFriendRequestsModel struct {
		*defaultFriendRequestsModel
	}
)

// FindByReqUidAndUserId 通过 reqUid（发起人ID） 和 userId（要添加的好友 ID） 查询好友申请的请求是否存在
func (c *customFriendRequestsModel) FindByReqUidAndUserId(ctx context.Context, reqUid string, userId string) (*FriendRequests, error) {

	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `user_id` = ?", friendRequestsRows, c.table)
	var friendRequests FriendRequests
	err := c.QueryRowNoCacheCtx(ctx, &friendRequests, query, reqUid, userId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &friendRequests, nil
}

// NewFriendRequestsModel returns a model for the database table.
func NewFriendRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestsModel {
	return &customFriendRequestsModel{
		defaultFriendRequestsModel: newFriendRequestsModel(conn, c, opts...),
	}
}
