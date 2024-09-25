package models

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	_                     UsersModel = (*customUsersModel)(nil)
	cacheUsersPhonePrefix            = "cache:users:id:"
)

type (

	/**
	要添加新方法则直接修改UserModel接口和给customUserModel添加方法
	*/
	// UserModel的接口继承了userModel_gen.go文件中的userModel接口
	UsersModel interface {
		usersModel
	}

	// customUsersModel继承了默认模型defaultUsersModel的所有功能，可以使用 defaultUserModel 的所有方法，并可以扩展更多自定义方法
	customUsersModel struct {
		*defaultUsersModel
	}
)

// FindByPhone 通过手机号查询用户详细信息
func (c *customUsersModel) FindByPhone(ctx context.Context, phone string) (*Users, error) {
	usersIdKey := fmt.Sprintf("%s%v", cacheUsersPhonePrefix, phone)
	var resp Users
	err := c.QueryRowCtx(ctx, &resp, usersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", usersRows, c.table)
		return conn.QueryRowCtx(ctx, v, query, phone)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}
