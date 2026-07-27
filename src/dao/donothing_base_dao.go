// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package dao

import (
	goctx "context"

	"github.com/beego/beego/v2/client/orm"
)

var _ BaseInterface = &DoNothingBase{}

type DoNothingBase struct {
}

func (d *DoNothingBase) UpdateWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (d *DoNothingBase) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) error {
	// TODO implement me
	panic("implement me")
}

func (d *DoNothingBase) InsertOrUpdateWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, md interface{}, colConflitAndArgs ...string) error {
	// TODO implement me
	panic("implement me")
}

func (d *DoNothingBase) InsertMultiWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, mds interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingBase) QueryMultiWithCtx(ctx goctx.Context, md interface{}, query string, args ...interface{}) error {
	return nil
}

func (d *DoNothingBase) Exec(ctx goctx.Context, query string, args ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingBase) ExecWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, query string, args ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingBase) DoTxWithCtx(ctx goctx.Context, task func(ctx goctx.Context, txOrm orm.TxOrmer) error) error {
	return task(ctx, nil)
}

func (d *DoNothingBase) List(md interface{}, opts ...QueryOption) error {
	return nil
}

func (d *DoNothingBase) Get(md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingBase) Delete(md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingBase) Insert(md interface{}) error {
	return nil
}

func (d *DoNothingBase) InsertWithOrm(ctx goctx.Context, proxy orm.QueryExecutor, md interface{}) error {
	return nil
}

func (d *DoNothingBase) Update(md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingBase) QueryOne(md interface{}, query string, args ...interface{}) error {
	return nil
}

func (d *DoNothingBase) QueryMulti(md interface{}, query string, args ...interface{}) error {
	return nil
}

func (d *DoNothingBase) InsertMulti(mds interface{}) error {
	return nil
}
