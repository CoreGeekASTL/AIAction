// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"GIDS/models/db"
)

/*
* 测试用例描述：TestWhiteListDao
* 预置条件：ormer 替换为 fakeOrm
* 操作步骤：
*     1. 分别调用 Count/GetByIMEIAndIMSI/InsertMulti/ClearAndInsert/ListAll
* 预期结果：
*     1. 均无错误返回
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestWhiteListDao(t *testing.T) {
	ormer = &fakeOrm{}
	d := NewWhiteListDao()

	count, err := d.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	err = d.GetByIMEIAndIMSI("625841245402541", "460011234567890")
	assert.NoError(t, err)

	list := []db.WhiteList{
		{Imei: "625841245402541", Imsi: "460011234567890", CreatedAt: "2026-08-11 00:00:00"},
	}
	err = d.InsertMulti(&list)
	assert.NoError(t, err)

	err = d.ClearAndInsert(&list)
	assert.NoError(t, err)

	var all []db.WhiteList
	err = d.ListAll(&all)
	assert.NoError(t, err)
}
