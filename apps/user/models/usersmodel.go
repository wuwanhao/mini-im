package models

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var (
	_                     UsersModel = (*customUsersModel)(nil)
	cacheUsersPhonePrefix            = "cache:users:phone:"
)

type (

	// UserModel的接口继承了userModel_gen.go文件中的userModel接口,要添加新方法则直接修改UserModel接口和给customUserModel添加方法
	UsersModel interface {
		usersModel
		FindByPhone(ctx context.Context, phone string) (*Users, error)
		FindUserByName(ctx context.Context, name string) ([]*Users, error)
		FindUserByIds(ctx context.Context, ids []string) ([]*Users, error)
	}

	// customUsersModel继承了默认模型defaultUsersModel的所有功能，可以使用 defaultUserModel 的所有方法，并可以扩展更多自定义方法
	customUsersModel struct {
		*defaultUsersModel
	}
)

// FindByPhone 通过手机号查询用户详细信息
func (c *customUsersModel) FindByPhone(ctx context.Context, phone string) (*Users, error) {
	usersPhoneKey := fmt.Sprintf("%s%v", cacheUsersPhonePrefix, phone)
	var resp Users
	err := c.QueryRowCtx(ctx, &resp, usersPhoneKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
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

// FindByName 通过用户名查询用户详细信息
func (c *customUsersModel) FindUserByName(ctx context.Context, name string) ([]*Users, error) {
	query := fmt.Sprintf("select %s from %s where `nickname` like ?", usersRows, c.table)
	var resp []*Users
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, fmt.Sprintf("%", name, "%"))
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}

}

// FindUserByIds 通过id查询用户详细信息
func (c *customUsersModel) FindUserByIds(ctx context.Context, ids []string) ([]*Users, error) {
	querySql := fmt.Sprintf("select %s from %s where `id` in ('%s')", usersRows, c.table, strings.Join(ids, ","))
	var resp []*Users
	return resp, c.QueryRowsNoCacheCtx(ctx, &resp, querySql)
}

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}
