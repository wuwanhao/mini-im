package models

import (
	"context"
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
	}

	customGroupsModel struct {
		*defaultGroupsModel
	}
)

// FindOneByName 根据群名查找群
func (c *customGroupsModel) FindOneByName(ctx context.Context, name string) (*Groups, error) {
	query := fmt.Sprintf("select %s from %s where `name` = ?", groupsFieldNames, c.table)
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

// NewGroupsModel returns a model for the database table.
func NewGroupsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupsModel {
	return &customGroupsModel{
		defaultGroupsModel: newGroupsModel(conn, c, opts...),
	}
}
