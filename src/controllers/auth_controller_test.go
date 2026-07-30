// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
* 测试用例描述：TestAuthControllerRouteInfo
* 预置条件：无
* 操作步骤：创建 AuthController 实例并调用 RouteInfo
* 预期结果：返回三个鉴权相关路由映射
 */
func TestAuthControllerRouteInfo(t *testing.T) {
	controller := &AuthController{}
	routeInfo := controller.RouteInfo()

	if reflect.ValueOf(routeInfo).IsZero() {
		t.Fatalf("RouteInfo is nil")
	}

	assert.Equal(t, "POST:AuthIMEI", routeInfo.RouteMapping["/auth/v1/authIMEI"])
	assert.Equal(t, "POST:ImportIMEIList", routeInfo.RouteMapping["/auth/v1/importIMEIList"])
	assert.Equal(t, "GET:ExportIMEIList", routeInfo.RouteMapping["/auth/v1/exportIMEIList"])
}
