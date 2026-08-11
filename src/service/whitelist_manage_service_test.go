// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package service
package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestManageService(d whiteListDaoInterface) *whiteListManageServiceImpl {
	return &whiteListManageServiceImpl{whiteListDao: d}
}

/*
* 测试用例描述：TestImportIMEIListFirstImport
* 预置条件：白名单表为空
* 操作步骤：
*     1. firstImport 导入 3 条合法记录
*     2. 表非空时再次 firstImport
* 预期结果：
*     1. 返回导入条数 3
*     2. 幂等拒绝，提示 not empty
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestImportIMEIListFirstImport(t *testing.T) {
	svc := newTestManageService(newFakeWhiteListDao())
	csvText := "625841245402541,460011234567890\n625841245402542,460011234567891\r\n625841245402543,460011234567892\n"
	count, err := svc.ImportIMEIList(strings.NewReader(csvText), operationFirstImport)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	_, err = svc.ImportIMEIList(strings.NewReader(csvText), operationFirstImport)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not empty")
}

/*
* 测试用例描述：TestImportIMEIListInvalidRow
* 预置条件：白名单管理服务实例
* 操作步骤：
*     1. 分别导入含非法 IMEI/IMSI、字段数错误、带 header 行的 CSV
* 预期结果：
*     1. 均整批拒绝返回错误，一条不入库
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestImportIMEIListInvalidRow(t *testing.T) {
	cases := []struct {
		name string
		csv  string
	}{
		{"imei 非15位纯数字", "6258412454025ABC,460011234567890\n"},
		{"imsi 非15位纯数字", "625841245402541,46001ABC4567890\n"},
		{"字段数错误", "625841245402541\n"},
		{"header行当数据校验失败", "IMEI,IMSI\n625841245402541,460011234567890\n"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			d := newFakeWhiteListDao()
			svc := newTestManageService(d)
			_, err := svc.ImportIMEIList(strings.NewReader(tt.csv), operationFirstImport)
			assert.Error(t, err)
			assert.Equal(t, 0, len(d.records))
		})
	}
}

/*
* 测试用例描述：TestImportIMEIListUpdateOverride
* 预置条件：白名单已有 3 条记录
* 操作步骤：
*     1. update 导入含空行的 2 条新记录
* 预期结果：
*     1. 空行跳过，旧记录被覆盖，表内仅剩 2 条新记录
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestImportIMEIListUpdateOverride(t *testing.T) {
	svc := newTestManageService(newFakeWhiteListDao())
	oldCSV := "625841245402541,460011234567890\n625841245402542,460011234567891\n625841245402543,460011234567892\n"
	_, err := svc.ImportIMEIList(strings.NewReader(oldCSV), operationFirstImport)
	assert.NoError(t, err)

	newCSV := "625841245402544,460011234567893\n\n625841245402545,460011234567894\n"
	count, err := svc.ImportIMEIList(strings.NewReader(newCSV), operationUpdate)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

/*
* 测试用例描述：TestExportIMEIList
* 预置条件：白名单含 2 条记录
* 操作步骤：
*     1. 调用导出接口
* 预期结果：
*     1. 返回无 header 的 IMEI/IMSI 两列 CSV 文本
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestExportIMEIList(t *testing.T) {
	svc := newTestManageService(newFakeWhiteListDao())
	csvText := "625841245402541,460011234567890\n625841245402542,460011234567891\n"
	_, err := svc.ImportIMEIList(strings.NewReader(csvText), operationFirstImport)
	assert.NoError(t, err)

	exported, err := svc.ExportIMEIList()
	assert.NoError(t, err)
	assert.Contains(t, exported, "625841245402541,460011234567890")
	assert.Contains(t, exported, "625841245402542,460011234567891")
	assert.NotContains(t, exported, "IMEI")
}
