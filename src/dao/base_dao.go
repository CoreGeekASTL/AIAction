/*
*  Copyright (c) Huawei Technologies Co., Ltd. 2025. All rights reserved.
 */

// Package dao 高斯数据库适配
package dao

import (
	goctx "context"
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

// defaultBulk 默认批量数
const defaultBulk = 100

var ormer orm.Ormer = &orm.DoNothingOrm{}

type Filter func() (string, []interface{})

type QueryOption struct {
	filters  []Filter
	orderBys []string
	limit    int
	offset   int
}

func NewQueryOption() *QueryOption {
	return &QueryOption{}
}

func (opt *QueryOption) Filter(key string, vals ...interface{}) *QueryOption {
	opt.filters = append(opt.filters, func() (string, []interface{}) {
		return key, vals
	})
	return opt
}

func (opt *QueryOption) OrderBy(o ...string) *QueryOption {
	opt.orderBys = o
	return opt
}

// Limit 设置分页
func (opt *QueryOption) Limit(limit, offset int, orderBy string) *QueryOption {
	if limit < 0 || offset < 0 {
		return opt
	}
	opt.OrderBy(orderBy)
	opt.limit = limit
	opt.offset = offset
	return opt
}

type BaseInterface interface {
	List(md interface{}, opts ...QueryOption) error
	Get(md interface{}, cols ...string) error
	Delete(md interface{}, cols ...string) error
	Insert(md interface{}) error
	InsertWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, md interface{}) error
	Update(md interface{}, cols ...string) error
	UpdateWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, md interface{}, cols ...string) error
	QueryOne(md interface{}, query string, args ...interface{}) error
	QueryMulti(md interface{}, query string, args ...interface{}) error
	QueryMultiWithCtx(ctx goctx.Context, md interface{}, query string, args ...interface{}) error
	Exec(ctx goctx.Context, query string, args ...interface{}) (int64, error)
	ExecWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, query string, args ...interface{}) (int64, error)
	InsertMulti(mds interface{}) error
	InsertMultiWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, mds interface{}) (int64, error)

	DoTxWithCtx(ctx goctx.Context, task func(ctx goctx.Context, txOrm orm.TxOrmer) error) error
}

var _ BaseInterface = &BaseDao{}

// BaseDao 数据库操作基础类
type BaseDao struct {
	EntityType interface{} // 数据表对应的entity类型
}

func (base *BaseDao) DoTxWithCtx(ctx goctx.Context, task func(ctx goctx.Context, txOrm orm.TxOrmer) error) error {
	return ormer.DoTxWithCtx(ctx, task)
}

func (base *BaseDao) UpdateWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, md interface{}, cols ...string) error {
	_, err := proxy.UpdateWithCtx(ctx, md, cols...)
	return err
}

// List 查询数据列表
func (base *BaseDao) List(md interface{}, opts ...QueryOption) error {
	qs := ormer.QueryTable(base.EntityType)
	if len(opts) > 0 {
		for _, item := range opts[0].filters {
			s, cols := item()
			qs = qs.Filter(s, cols...)
		}

		if len(opts[0].orderBys) > 0 {
			qs = qs.OrderBy(opts[0].orderBys...)
		}

		if opts[0].limit > 0 {
			qs = qs.Limit(opts[0].limit, opts[0].offset)
		}
	}

	_, err := qs.All(md)
	return err
}

// Get 查询数据表
func (base *BaseDao) Get(md interface{}, cols ...string) error {
	return ormer.Read(md, cols...)
}

// Delete 删除数据表
func (base *BaseDao) Delete(md interface{}, cols ...string) error {
	_, err := ormer.Delete(md, cols...)
	return err
}

// Insert 新增数据表
func (base *BaseDao) Insert(md interface{}) error {
	return base.InsertWithOrm(goctx.Background(), ormer, md)
}

func (base *BaseDao) InsertWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, md interface{}) error {
	_, err := proxy.InsertWithCtx(ctx, md)
	if errors.Is(err, orm.ErrLastInsertIdUnavailable) {
		return nil
	}
	return err
}

// Update 更新数据表
func (base *BaseDao) Update(md interface{}, cols ...string) error {
	return base.UpdateWithOrm(goctx.Background(), ormer, md, cols...)
}

func (base *BaseDao) QueryOne(md interface{}, query string, args ...interface{}) error {
	return ormer.Raw(query, args...).QueryRow(md)
}

func (base *BaseDao) QueryMulti(md interface{}, query string, args ...interface{}) error {
	return base.QueryMultiWithCtx(goctx.Background(), md, query, args...)
}

func (base *BaseDao) QueryMultiWithCtx(ctx goctx.Context, md interface{}, query string, args ...interface{}) error {
	_, err := ormer.RawWithCtx(ctx, query, args...).QueryRows(md)
	return err
}

func (base *BaseDao) Exec(ctx goctx.Context, query string, args ...interface{}) (int64, error) {
	return base.ExecWithOrm(ctx, ormer, query, args...)
}

func (base *BaseDao) ExecWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, query string, args ...interface{}) (int64, error) {
	ret, err := proxy.RawWithCtx(ctx, query, args...).Exec()
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

func (base *BaseDao) InsertMulti(mds interface{}) error {
	_, err := base.InsertMultiWithOrm(goctx.Background(), ormer, mds)
	return err
}

func (base *BaseDao) InsertMultiWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, mds interface{}) (int64, error) {
	return proxy.InsertMultiWithCtx(ctx, defaultBulk, mds)
}
