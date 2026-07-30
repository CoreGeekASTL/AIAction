// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package service

import (
	"errors"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	"github.com/stretchr/testify/assert"

	"GIDS/models/db"
)

// fakeWhiteListStore 白名单存取 fake，记录调用次数验证缓存命中
type fakeWhiteListStore struct {
	records     map[string]db.AuthWhitelist
	countErr    error
	getErr      error
	countCalls  int
	getCalls    int
	insertCalls int
	clearCalls  int
	listCalls   int
}

func newFakeWhiteListStore() *fakeWhiteListStore {
	return &fakeWhiteListStore{records: make(map[string]db.AuthWhitelist)}
}

func (f *fakeWhiteListStore) Count() (int64, error) {
	f.countCalls++
	return int64(len(f.records)), f.countErr
}

func (f *fakeWhiteListStore) GetByIMEI(imei string) (*db.AuthWhitelist, error) {
	f.getCalls++
	if f.getErr != nil {
		return nil, f.getErr
	}
	wl, ok := f.records[imei]
	if !ok {
		return nil, orm.ErrNoRows
	}
	return &wl, nil
}

func (f *fakeWhiteListStore) InsertMulti(records []db.AuthWhitelist) error {
	f.insertCalls++
	for _, r := range records {
		f.records[r.IMEI] = r
	}
	return nil
}

func (f *fakeWhiteListStore) ClearAndInsert(records []db.AuthWhitelist) error {
	f.clearCalls++
	f.records = make(map[string]db.AuthWhitelist)
	for _, r := range records {
		f.records[r.IMEI] = r
	}
	return nil
}

func (f *fakeWhiteListStore) ListAll() ([]db.AuthWhitelist, error) {
	f.listCalls++
	list := make([]db.AuthWhitelist, 0, len(f.records))
	for _, r := range f.records {
		list = append(list, r)
	}
	return list, nil
}

func newTestAuthService(store whiteListDao) *authServiceImpl {
	return &authServiceImpl{store: store, cache: newAuthCache()}
}

const (
	testIMEI  = "625841245402541"
	testIMSI  = "460011234567890"
	testIMEI2 = "625841245402599"
	testIMSI2 = "460019999999999"
)

/*
* 测试用例描述：TestAuthCheck 鉴权各分支
* 预置条件：fake 白名单存取
* 操作步骤：分别构造格式非法/逃生态/命中/IMSI不符/IMEI未命中场景调用 Check
* 预期结果：格式非法返回(false,false)；逃生态与组合命中放行；其余拒绝
 */
func TestAuthCheck(t *testing.T) {
	store := newFakeWhiteListStore()
	svc := newTestAuthService(store)

	pass, valid := svc.Check("123", testIMSI)
	assert.False(t, pass)
	assert.False(t, valid)

	pass, valid = svc.Check(testIMEI, testIMSI)
	assert.True(t, pass)
	assert.True(t, valid)

	store.records[testIMEI] = db.AuthWhitelist{IMEI: testIMEI, IMSI: testIMSI}
	svc2 := newTestAuthService(store)

	pass, _ = svc2.Check(testIMEI, testIMSI)
	assert.True(t, pass)

	pass, _ = svc2.Check(testIMEI, testIMSI2)
	assert.False(t, pass)

	pass, _ = svc2.Check(testIMEI2, testIMSI)
	assert.False(t, pass)
}

/*
* 测试用例描述：TestAuthCheckCache 缓存命中不再回源、ClearCache 后重新回源
* 预置条件：白名单含一条记录
* 操作步骤：同一组合连续两次 Check，随后 ClearCache 再 Check
* 预期结果：第二次无新增 DB 调用；ClearCache 后回源次数增加
 */
func TestAuthCheckCache(t *testing.T) {
	store := newFakeWhiteListStore()
	store.records[testIMEI] = db.AuthWhitelist{IMEI: testIMEI, IMSI: testIMSI}
	svc := newTestAuthService(store)

	svc.Check(testIMEI, testIMSI)
	countCalls, getCalls := store.countCalls, store.getCalls

	pass, _ := svc.Check(testIMEI, testIMSI)
	assert.True(t, pass)
	assert.Equal(t, countCalls, store.countCalls)
	assert.Equal(t, getCalls, store.getCalls)

	pass, _ = svc.Check(testIMEI2, testIMSI2)
	assert.False(t, pass)
	getCalls = store.getCalls
	pass, _ = svc.Check(testIMEI2, testIMSI2)
	assert.False(t, pass)
	assert.Equal(t, getCalls, store.getCalls)

	svc.ClearCache()
	svc.Check(testIMEI, testIMSI)
	assert.Greater(t, store.countCalls, countCalls)
}

/*
* 测试用例描述：TestAuthCheckStoreError DB 异常 fail-open
* 预置条件：store Count/GetByIMEI 返回错误
* 操作步骤：调用 Check
* 预期结果：放行，避免 DB 故障阻断主流程
 */
func TestAuthCheckStoreError(t *testing.T) {
	store := newFakeWhiteListStore()
	store.countErr = errors.New("db down")
	svc := newTestAuthService(store)
	pass, _ := svc.Check(testIMEI, testIMSI)
	assert.True(t, pass)

	store.countErr = nil
	store.records[testIMEI] = db.AuthWhitelist{IMEI: testIMEI, IMSI: testIMSI}
	store.getErr = errors.New("db down")
	svc2 := newTestAuthService(store)
	pass, _ = svc2.Check(testIMEI, testIMSI)
	assert.True(t, pass)
}

/*
* 测试用例描述：TestAuthCacheEvict 容量超限惰性清理最旧 500 条
* 预置条件：空缓存
* 操作步骤：写入容量上限+1 条
* 预期结果：清理后剩余 501 条，最早写入的条目被淘汰
 */
func TestAuthCacheEvict(t *testing.T) {
	cache := newAuthCache()
	for i := 0; i < authCacheCapacity+1; i++ {
		key := string(rune('a'+i%26)) + string(rune('A'+i/26)) + "_key"
		cache.set(key, true)
	}
	assert.Equal(t, authCacheCapacity+1-authCacheCleanCount, len(cache.items))
}
