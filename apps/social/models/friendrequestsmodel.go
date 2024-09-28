package models

import (
	"app/pkg/constants"
	"context"
	"fmt"
	"github.com/pkg/errors"
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
		ListNoHandler(ctx context.Context, userId string) ([]*FriendRequests, error)
		Trans(ctx context.Context, fn func(ctx context.Context, conn sqlx.SqlConn)) error
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

// ListNoHandler 获取当前用户未处理的好友申请列表
func (c *customFriendRequestsModel) ListNoHandler(ctx context.Context, userId string) ([]*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `handle_result` = ?", friendRequestsRows, c.table)
	var list []*FriendRequests
	err := c.QueryRowsNoCacheCtx(ctx, &list, query, userId, constants.NoHandlerResult)
	if err != nil {
		return nil, errors.Wrapf(err, "ListNoHandler error, userId: %v", userId)
	}
	return list, nil
}

func (c *customFriendRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, conn sqlx.SqlConn)) error {
	return c.TransactCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) error {
		return fn(ctx, conn)
	})
} // todo

// NewFriendRequestsModel returns a model for the database table.
func NewFriendRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestsModel {
	return &customFriendRequestsModel{
		defaultFriendRequestsModel: newFriendRequestsModel(conn, c, opts...),
	}
}
