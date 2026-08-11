// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package service
package service

import (
	"errors"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	"github.com/stretchr/testify/assert"

	"GIDS/models/db"
)

// fakeWhiteListDao 白名单 DAO 假实现，内存态模拟表数据
type fakeWhiteListDao struct {
	records map[string]db.WhiteList
	countErr error
}

func newFakeWhiteListDao() *fakeWhiteListDao {
	return &fakeWhiteListDao{records: make(map[string]db.WhiteList)}
}

func (f *fakeWhiteListDao) Count() (int64, error) {
	if f.countErr != nil {
		return 0, f.countErr
	}
	return int64(len(f.records)), nil
}

func (f *fakeWhiteListDao) GetByIMEIAndIMSI(imei, imsi string) error {
	if w, ok := f.records[imei]; ok && w.Imsi == imsi {
		return nil
	}
	return orm.ErrNoRows
}

func (f *fakeWhiteListDao) InsertMulti(list *[]db.WhiteList) error {
	for _, w := range *list {
		f.records[w.Imei] = w
	}
	return nil
}

func (f *fakeWhiteListDao) ClearAndInsert(list *[]db.WhiteList) error {
	f.records = make(map[string]db.WhiteList)
	return f.InsertMulti(list)
}

func (f *fakeWhiteListDao) ListAll(list *[]db.WhiteList) error {
	for _, w := range f.records {
		*list = append(*list, w)
	}
	return nil
}

func newTestAuthService(d whiteListDaoInterface) *authServiceImpl {
	return &authServiceImpl{
		whiteListDao: d,
		cache:        &authCache{entries: make(map[string]cacheEntry)},
	}
}

/*
* 测试用例描述：TestAuthIMEIFormatInvalid
* 预置条件：鉴权服务实例
* 操作步骤：
*     1. 分别传入 14 位、16 位、含字母、空串的 IMEI/IMSI
* 预期结果：
*     1. 均短路返回 allowed=false, formatValid=false
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthIMEIFormatInvalid(t *testing.T) {
	svc := newTestAuthService(newFakeWhiteListDao())
	cases := []struct {
		name string
		imei string
		imsi string
	}{
		{"imei 14位", "62584124540254", "460011234567890"},
		{"imei 16位", "6258412454025411", "460011234567890"},
		{"imei 含字母", "6258412454025AB", "460011234567890"},
		{"imei 空串", "", "460011234567890"},
		{"imsi 14位", "625841245402541", "46001123456789"},
		{"imsi 16位", "625841245402541", "4600112345678901"},
		{"imsi 含字母", "625841245402541", "46001ABC4567890"},
		{"imsi 空串", "625841245402541", ""},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			allowed, formatValid := svc.AuthIMEI(tt.imei, tt.imsi)
			assert.False(t, allowed)
			assert.False(t, formatValid)
		})
	}
}

/*
* 测试用例描述：TestAuthIMEIEscapeMode
* 预置条件：白名单表为空（逃生态）
* 操作步骤：
*     1. 合法 15 位 IMEI/IMSI 鉴权
*     2. 相同参数再次鉴权
* 预期结果：
*     1. 放行
*     2. 命中缓存放行标记仍放行
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthIMEIEscapeMode(t *testing.T) {
	svc := newTestAuthService(newFakeWhiteListDao())
	imei := "625841245402541"
	imsi := "460011234567890"

	allowed, formatValid := svc.AuthIMEI(imei, imsi)
	assert.True(t, allowed)
	assert.True(t, formatValid)

	allowed, formatValid = svc.AuthIMEI(imei, imsi)
	assert.True(t, allowed)
	assert.True(t, formatValid)
}

/*
* 测试用例描述：TestAuthIMEIJointMatch
* 预置条件：白名单含一条 IMEI+IMSI 组合
* 操作步骤：
*     1. 联合命中鉴权
*     2. IMEI 命中 IMSI 未命中
*     3. 重复未命中鉴权
* 预期结果：
*     1. 放行
*     2. 拒绝并缓存 false
*     3. 缓存 false 直接返回拒绝
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthIMEIJointMatch(t *testing.T) {
	d := newFakeWhiteListDao()
	d.records["625841245402541"] = db.WhiteList{Imei: "625841245402541", Imsi: "460011234567890"}
	svc := newTestAuthService(d)

	allowed, _ := svc.AuthIMEI("625841245402541", "460011234567890")
	assert.True(t, allowed)

	allowed, formatValid := svc.AuthIMEI("625841245402541", "460019999999999")
	assert.False(t, allowed)
	assert.True(t, formatValid)

	allowed, _ = svc.AuthIMEI("625841245402541", "460019999999999")
	assert.False(t, allowed)
}

/*
* 测试用例描述：TestAuthIMEIDBErrorSafeReject
* 预置条件：白名单 DAO Count 返回异常
* 操作步骤：
*     1. 合法 IMEI/IMSI 鉴权
* 预期结果：
*     1. 缓存未命中且 DB 异常，安全优先拒绝 allowed=false, formatValid=true
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthIMEIDBErrorSafeReject(t *testing.T) {
	d := newFakeWhiteListDao()
	d.countErr = errors.New("db connection failed")
	svc := newTestAuthService(d)

	allowed, formatValid := svc.AuthIMEI("625841245402541", "460011234567890")
	assert.False(t, allowed)
	assert.True(t, formatValid)
}
