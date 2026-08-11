// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package dao

import (
	goctx "context"

	"github.com/beego/beego/v2/client/orm"

	"GIDS/models/db"
)

// WhiteListDao 白名单数据访问对象
type WhiteListDao struct {
	BaseInterface
}

// NewWhiteListDao 构造白名单 DAO
func NewWhiteListDao() *WhiteListDao {
	dao := &WhiteListDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.WhiteList{},
	}
	return dao
}

// Count 统计白名单记录数
func (d *WhiteListDao) Count() (int64, error) {
	return ormer.QueryTable(&db.WhiteList{}).Count()
}

// GetByIMEIAndIMSI IMEI+IMSI 联合精确匹配查询，未命中返回 orm.ErrNoRows
func (d *WhiteListDao) GetByIMEIAndIMSI(imei, imsi string) error {
	var w db.WhiteList
	return ormer.QueryTable(&db.WhiteList{}).Filter("Imei", imei).Filter("Imsi", imsi).One(&w)
}

// InsertMulti 批量插入白名单记录
func (d *WhiteListDao) InsertMulti(list *[]db.WhiteList) error {
	return d.BaseInterface.InsertMulti(list)
}

// ClearAndInsert 事务内清表并批量插入，失败整批回滚
func (d *WhiteListDao) ClearAndInsert(list *[]db.WhiteList) error {
	return d.DoTxWithCtx(goctx.TODO(), func(ctx goctx.Context, txOrm orm.TxOrmer) error {
		if _, err := txOrm.Raw("DELETE FROM t_white_list").Exec(); err != nil {
			return err
		}
		if len(*list) == 0 {
			return nil
		}
		_, err := d.InsertMultiWithOrm(ctx, txOrm, list)
		return err
	})
}

// ListAll 全量查询白名单记录
func (d *WhiteListDao) ListAll(list *[]db.WhiteList) error {
	return d.List(list)
}
