/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package event

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"GIDS/models/events"
)

/*
* 测试用例描述：测试 InitFactory 函数
* 预置条件：工厂实例未初始化
* 操作步骤：
*     1. 调用 InitFactory 函数
* 预期结果：
*     1. 工厂实例被正确初始化，factory 字段为空
* 修改历史：
*     1. 2025-12-24 新建测试用例
 */
func TestInitFactory(t *testing.T) {
	InitFactory()
	assert.NotNil(t, factoryInstance)
	assert.Equal(t, 0, len(factoryInstance.(*eventStorageFactory).factory))
}

/*
* 测试用例描述：测试 GetFactory 函数
* 预置条件：工厂实例已初始化
* 操作步骤：
*     1. 调用 GetFactory 函数
* 预期结果：
*     1. 返回的工厂实例与全局工厂实例一致
* 修改历史：
*     1. 2025-12-24 新建测试用例
 */
func TestGetFactory(t *testing.T) {
	factory := GetFactory()
	assert.NotNil(t, factory)
	assert.Equal(t, factoryInstance, factory)
}

/*
* 测试用例描述：测试 Register 函数
* 预置条件：工厂实例已初始化，存储实例已创建
* 操作步骤：
*     1. 调用 Register 函数，注册存储实例
* 预期结果：
*     1. 存储实例被正确注册到工厂中
* 修改历史：
*     1. 2025-12-24 新建测试用例
 */
func TestRegister(t *testing.T) {
	testStorage := &mockStorage{}

	factory, ok := GetFactory().(*eventStorageFactory)
	if !ok {
		t.Fatalf("type assertion failed")
	}
	factory.Register("test", testStorage)

	st, err := factory.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, testStorage, st)
}

/*
* 测试用例描述：测试 Get 函数
* 预置条件：工厂实例已初始化，存储实例已注册
* 操作步骤：
*     1. 调用 Get 函数，获取存储实例
* 预期结果：
*     1. 返回的存储实例与注册的实例一致
* 修改历史：
*     1. 2025-12-24 新建测试用例
 */
func TestGet(t *testing.T) {
	testStorage := &mockStorage{}

	factory, ok := GetFactory().(*eventStorageFactory)
	if !ok {
		t.Fatalf("type assertion failed")
	}
	factory.Register("test", testStorage)

	// 测试获取存储
	st, err := factory.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, testStorage, st)

	// 测试获取不存在的存储
	_, err = factory.Get("non-existent")
	assert.Error(t, err)
	assert.Equal(t, "eventStorage non-existent not found in factory", err.Error())

	// 测试获取存储为 nil 的情况
	factory.factory["nil"] = nil
	_, err = factory.Get("nil")
	assert.Error(t, err)
	assert.Equal(t, "nil eventStorage nil found in factory", err.Error())
}

type mockStorage struct{}

func (ts *mockStorage) Record(logInfo *events.Info) error {
	return nil
}
