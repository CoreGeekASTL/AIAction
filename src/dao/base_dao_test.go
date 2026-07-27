/*
 *  Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved.
 */

package dao

import (
	goctx "context"
	"database/sql"
	"database/sql/driver"
	"testing"

	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/clauses/order_clause"
	convey "github.com/smartystreets/goconvey/convey"

	models "GIDS/models/db"
	testutil "GIDS/test/util"
)

var _ orm.QuerySeter = &DoNothingQuerySeter{}

type DoNothingQuerySeter struct {
}

func (d *DoNothingQuerySeter) Filter(s string, i ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) FilterRaw(s string, s2 string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) Exclude(s string, i ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) SetCond(condition *orm.Condition) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) GetCond() *orm.Condition {
	return nil
}

func (d *DoNothingQuerySeter) Limit(limit interface{}, args ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) Offset(offset interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) GroupBy(exprs ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) OrderBy(exprs ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) OrderClauses(orders ...*order_clause.Order) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) ForceIndex(indexes ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) UseIndex(indexes ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) IgnoreIndex(indexes ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) RelatedSel(params ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) Distinct() orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) ForUpdate() orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySeter) Count() (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) CountWithCtx(ctx goctx.Context) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) Exist() bool {
	return true
}

func (d *DoNothingQuerySeter) ExistWithCtx(ctx goctx.Context) bool {
	return true
}

func (d *DoNothingQuerySeter) Update(values orm.Params) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) UpdateWithCtx(ctx goctx.Context, values orm.Params) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) Delete() (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) DeleteWithCtx(ctx goctx.Context) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) PrepareInsert() (orm.Inserter, error) {
	return nil, nil
}

func (d *DoNothingQuerySeter) PrepareInsertWithCtx(ctx goctx.Context) (orm.Inserter, error) {
	return nil, nil
}

func (d *DoNothingQuerySeter) All(container interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) AllWithCtx(ctx goctx.Context, container interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) One(container interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingQuerySeter) OneWithCtx(ctx goctx.Context, container interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingQuerySeter) Values(results *[]orm.Params, exprs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) ValuesWithCtx(ctx goctx.Context, results *[]orm.Params, exprs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) ValuesList(results *[]orm.ParamsList, exprs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) ValuesListWithCtx(ctx goctx.Context, results *[]orm.ParamsList, exprs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) ValuesFlat(result *orm.ParamsList, expr string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) ValuesFlatWithCtx(ctx goctx.Context, result *orm.ParamsList, expr string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) RowsToMap(result *orm.Params, keyCol, valueCol string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySeter) Aggregate(s string) orm.QuerySeter {
	return d
}

var _ orm.RawSeter = &DoNothingRawSeter{}

type DoNothingRawSeter struct {
}

func (d *DoNothingRawSeter) Exec() (sql.Result, error) {
	return driver.RowsAffected(0), nil
}

func (d *DoNothingRawSeter) QueryRow(containers ...interface{}) error {
	return nil
}

func (d *DoNothingRawSeter) QueryRows(containers ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSeter) SetArgs(...interface{}) orm.RawSeter {
	return d
}

func (d *DoNothingRawSeter) Values(container *[]orm.Params, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSeter) ValuesList(container *[]orm.ParamsList, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSeter) ValuesFlat(container *orm.ParamsList, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSeter) RowsToMap(result *orm.Params, keyCol, valueCol string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSeter) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSeter) Prepare() (orm.RawPreparer, error) {
	return nil, nil
}

type fakeOrm struct {
	orm.DoNothingOrm
}

func (f *fakeOrm) QueryTable(ptrStructOrTableName interface{}) orm.QuerySeter {
	return &DoNothingQuerySeter{}
}

func (f *fakeOrm) Raw(query string, args ...interface{}) orm.RawSeter {
	return &DoNothingRawSeter{}
}

func (f *fakeOrm) RawWithCtx(ctx goctx.Context, query string, args ...interface{}) orm.RawSeter {
	return &DoNothingRawSeter{}
}

func TestBaseDao(t *testing.T) {
	convey.Convey("test base dao", t, func() {
		ormer = &fakeOrm{}
		var b BaseInterface

		convey.Convey("test base dao", func() {
			b = &BaseDao{
				EntityType: &models.RouterAPPConfig{},
			}
		})
		convey.Convey("test mock dao", func() {
			b = &DoNothingBase{}
		})
		convey.It("test base dao list", func() {
			var (
				md      []models.RouterAPPConfig
				orderBy = ""
			)
			err := b.List(&md, *testutil.NewQueryOption().Filter("", nil).OrderBy(orderBy))
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao get", func() {
			var md models.RouterAPPConfig
			err := b.Get(&md)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao insert", func() {
			var md models.RouterAPPConfig
			err := b.Insert(&md)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao update", func() {
			var md models.RouterAPPConfig
			err := b.Update(&md)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao delete", func() {
			var md models.RouterAPPConfig
			err := b.Delete(&md)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao query", func() {
			var md models.RouterAPPConfig
			err := b.QueryOne(&md, "select * from t_route_app_configs where id = 1")
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao query multi", func() {
			var md []models.RouterAPPConfig
			err := b.QueryMulti(&md, "select * from t_route_app_configs")
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao exec", func() {
			_, err := b.Exec(goctx.Background(), "delete from t_route_app_configs where id = 1")
			convey.So(err, convey.ShouldBeNil)
		})
		convey.It("test base dao insert multi", func() {
			var md []models.RouterAPPConfig
			err := b.InsertMulti(&md)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.It("test base dao do with ctx", func() {
			err := b.DoTxWithCtx(goctx.Background(), func(ctx goctx.Context, txOrm orm.TxOrmer) error {
				return nil
			})
			convey.So(err, convey.ShouldBeNil)
		})

		convey.It("test base dao insert with ctx", func() {
			var md models.RouterAPPConfig
			err := b.InsertWithOrm(goctx.Background(), ormer, &md)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
