// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
* 测试用例描述： TestEventController_RouteInfo
* 预置条件：无
* 操作步骤：
*     1. 创建 EventController 实例
*     2. 调用 RouteInfo 方法
*     3. 检查 RouteInfo 返回的结构体是否为空
* 预期结果：
*     1. 返回正确的 RouteInfo 结构
* 修改历史：
*     1. 2025-8-28 新建测试用例
 */
func TestEventControllerRouteInfo(t *testing.T) {
	controller := &EventController{}
	routeInfo := controller.RouteInfo()

	// 检查 routeInfo 是否为空，避免空指针解引用
	if reflect.ValueOf(routeInfo).IsZero() {
		t.Fatalf("RouteInfo is nil")
	}

	assert.Equal(t, "POST:SendClientEvent", routeInfo.RouteMapping["/app-api/center/public/client/sendClientEvent"])
	assert.Equal(t, "POST:SendAppUseTimesEvent", routeInfo.RouteMapping["/app-api/center/public/client/sendAppUseTimesEvent"])
}
