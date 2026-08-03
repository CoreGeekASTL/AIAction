// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"GIDS/common/constants/retcode"
	"GIDS/models/db"
)

func newTestManageService(store whiteListDao) *authManageServiceImpl {
	return &authManageServiceImpl{store: store, auth: newTestAuthService(store)}
}

const testCSV = "625841245402541,460011234567890\n625841245402542,460011234567891\n"

/*
* 测试用例描述：TestImportFirstImport 首次导入成功、重复首次导入拒绝
* 预置条件：白名单表为空
* 操作步骤：firstImport 导入合法 CSV，再次 firstImport
* 预期结果：首次 code=200 条数=2；重复 code=-1 且 msg 含 not empty
 */
func TestImportFirstImport(t *testing.T) {
	store := newFakeWhiteListStore()
	svc := newTestManageService(store)

	count, code, _ := svc.Import(strings.NewReader(testCSV), operationFirstImport)
	assert.Equal(t, retcode.Success, code)
	assert.Equal(t, 2, count)

	_, code, msg := svc.Import(strings.NewReader(testCSV), operationFirstImport)
	assert.Equal(t, retcode.InternalFailed, code)
	assert.Contains(t, strings.ToLower(msg), "not empty")
}

/*
* 测试用例描述：TestImportUpdate 覆盖导入清表重写
* 预置条件：表内已有旧记录
* 操作步骤：update 导入新 CSV
* 预期结果：code=200，旧记录被覆盖
 */
func TestImportUpdate(t *testing.T) {
	store := newFakeWhiteListStore()
	store.records[testIMEI2] = db.AuthWhitelist{IMEI: testIMEI2, IMSI: testIMSI2}
	svc := newTestManageService(store)

	count, code, _ := svc.Import(strings.NewReader(testCSV), operationUpdate)
	assert.Equal(t, retcode.Success, code)
	assert.Equal(t, 2, count)
	assert.Equal(t, 1, store.clearCalls)
	_, ok := store.records[testIMEI2]
	assert.False(t, ok)
}

/*
* 测试用例描述：TestImportInvalidParam 参数与内容校验失败分支
* 预置条件：白名单表为空
* 操作步骤：分别使用非法 operation、非法 IMEI、非法 IMSI、文件内重复组合导入
* 预期结果：operation 非法 code=-2；内容非法 code=-1 且整批不加载
 */
func TestImportInvalidParam(t *testing.T) {
	store := newFakeWhiteListStore()
	svc := newTestManageService(store)

	_, code, _ := svc.Import(strings.NewReader(testCSV), "badOp")
	assert.Equal(t, retcode.ClientFailed, code)

	badCases := []string{
		"62584124540254,460011234567890\n",
		"6258412454025411,460011234567890\n",
		"6258412454025ABC,460011234567890\n",
		"625841245402541,46001123456789\n",
		"625841245402541,46001ABC4567890\n",
		"625841245402541,460011234567890\n625841245402541,460011234567890\n",
	}
	for _, bad := range badCases {
		_, code, _ = svc.Import(strings.NewReader(bad), operationFirstImport)
		assert.Equal(t, retcode.InternalFailed, code, "csv: %q", bad)
	}
	assert.Equal(t, 0, store.insertCalls)
}

/*
* 测试用例描述：TestImportClearCache 导入成功后鉴权缓存被清空
* 预置条件：鉴权缓存中已有逃生态放行标记
* 操作步骤：firstImport 导入白名单后再鉴权未命中组合
* 预期结果：导入后缓存失效，未命中组合被拒绝
 */
func TestImportClearCache(t *testing.T) {
	store := newFakeWhiteListStore()
	auth := newTestAuthService(store)
	svc := &authManageServiceImpl{store: store, auth: auth}

	isPass, _ := auth.Check(testIMEI2, testIMSI2)
	assert.True(t, isPass)

	_, code, _ := svc.Import(strings.NewReader(testCSV), operationFirstImport)
	assert.Equal(t, retcode.Success, code)

	isPass, _ = auth.Check(testIMEI2, testIMSI2)
	assert.False(t, isPass)
	isPass, _ = auth.Check(testIMEI, testIMSI)
	assert.True(t, isPass)
}

/*
* 测试用例描述：TestExport 导出生成带表头的全量 CSV
* 预置条件：表内两条记录
* 操作步骤：调用 Export
* 预期结果：首行为 IMEI,IMSI 表头，数据行含两条记录
 */
func TestExport(t *testing.T) {
	store := newFakeWhiteListStore()
	store.records[testIMEI] = db.AuthWhitelist{IMEI: testIMEI, IMSI: testIMSI}
	store.records[testIMEI2] = db.AuthWhitelist{IMEI: testIMEI2, IMSI: testIMSI2}
	svc := newTestManageService(store)

	text, err := svc.Export()
	assert.Nil(t, err)
	lines := strings.Split(strings.TrimSpace(text), "\n")
	assert.Equal(t, "IMEI,IMSI", lines[0])
	assert.Equal(t, 3, len(lines))
	assert.Contains(t, text, testIMEI+","+testIMSI)
}
