// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"GIDS/models/req"
)

/*
* 测试用例描述：TestAuthControllerRouteInfo
* 预置条件：无
* 操作步骤：
*     1. 创建 AuthController 实例并调用 RouteInfo
* 预期结果：
*     1. 返回三个接口的路由映射
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthControllerRouteInfo(t *testing.T) {
	controller := &AuthController{}
	routeInfo := controller.RouteInfo()

	assert.Equal(t, "POST:AuthIMEI", routeInfo.RouteMapping["/auth/v1/authIMEI"])
	assert.Equal(t, "POST:ImportIMEIList", routeInfo.RouteMapping["/auth/v1/importIMEIList"])
	assert.Equal(t, "GET:ExportIMEIList", routeInfo.RouteMapping["/auth/v1/exportIMEIList"])
}

/*
* 测试用例描述：TestAuthIMEIRequestValidate
* 预置条件：无
* 操作步骤：
*     1. 分别构造 IMEI/IMSI 缺失与完整的请求并调用 Validate
* 预期结果：
*     1. 缺失报错，完整通过
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthIMEIRequestValidate(t *testing.T) {
	assert.Error(t, req.AuthIMEIRequest{IMEI: "", IMSI: "460011234567890"}.Validate())
	assert.Error(t, req.AuthIMEIRequest{IMEI: "625841245402541", IMSI: ""}.Validate())
	assert.NoError(t, req.AuthIMEIRequest{IMEI: "625841245402541", IMSI: "460011234567890"}.Validate())
}
