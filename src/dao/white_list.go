// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package dao

import (
	goctx "context"

	"github.com/beego/beego/v2/client/orm"

	"GIDS/models/db"
)

// WhiteListDao 白名单表读写，继承 BaseInterface
type WhiteListDao struct {
	BaseInterface
}

func NewWhiteListDao() *WhiteListDao {
	dao := &WhiteListDao{}
	dao.BaseInterface = &BaseDao{EntityType: &db.AuthWhitelist{}}
	return dao
}

// Count 统计白名单记录数，逃生态判定依据
func (d *WhiteListDao) Count() (int64, error) {
	var count int64
	err := d.QueryOne(&count, "SELECT COUNT(*) FROM t_white_list")
	return count, err
}

// GetByIMEI 按 IMEI 主键查询整行（含 IMSI），未命中返回 orm.ErrNoRows
func (d *WhiteListDao) GetByIMEI(imei string) (*db.AuthWhitelist, error) {
	record := &db.AuthWhitelist{IMEI: imei}
	if err := d.Get(record); err != nil {
		return nil, err
	}
	return record, nil
}

// InsertMulti 批量插入白名单记录
func (d *WhiteListDao) InsertMulti(records []db.AuthWhitelist) error {
	if len(records) == 0 {
		return nil
	}
	return d.BaseInterface.InsertMulti(records)
}

// ClearAndInsert 事务清表 + 批量插入，update 覆盖导入用，失败整体回滚
func (d *WhiteListDao) ClearAndInsert(records []db.AuthWhitelist) error {
	return d.DoTxWithCtx(goctx.Background(), func(ctx goctx.Context, txOrm orm.TxOrmer) error {
		if _, err := d.ExecWithOrm(ctx, txOrm, "DELETE FROM t_white_list"); err != nil {
			return err
		}
		if len(records) == 0 {
			return nil
		}
		_, err := d.InsertMultiWithOrm(ctx, txOrm, records)
		return err
	})
}

// ListAll 全量查询白名单，导出用
func (d *WhiteListDao) ListAll() ([]db.AuthWhitelist, error) {
	list := make([]db.AuthWhitelist, 0)
	if err := d.List(&list); err != nil {
		return nil, err
	}
	return list, nil
}
