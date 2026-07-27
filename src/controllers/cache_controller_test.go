// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
* 测试用例描述：TestCacheController_RouteInfo
* 预置条件：
* 操作步骤：
*     1. 创建 CacheController 实例
*     2. 调用 RouteInfo 方法
* 预期结果：
*     1. 返回正确的 RouteInfo 结构
* 修改历史：
*     1. 2025-8-27 新建测试用例
 */
func TestCacheControllerRouteInfo(t *testing.T) {
	controller := &CacheController{}
	routeInfo := controller.RouteInfo()

	assert.Equal(t, "POST:DeleteCache", routeInfo.RouteMapping["/app-api/devicetcp/cache/deleteCache"])
}
