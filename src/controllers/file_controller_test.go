// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
* 测试用例描述：TestFileController_RouteInfo
* 预置条件：
* 操作步骤：
*     1. 创建 FileController 实例
*     2. 调用 RouteInfo 方法
* 预期结果：
*     1. 返回正确的 RouteInfo 结构
* 修改历史：
*     1. 2025-8-26 新建测试用例
 */
func TestFileControllerRouteInfo(t *testing.T) {
	controller := &FileController{}
	routeInfo := controller.RouteInfo()

	assert.Equal(t, "POST:HandleUpload", routeInfo.RouteMapping["/app-api/control/file/upload"])
	assert.Equal(t, "GET:HandleDownload", routeInfo.RouteMapping["/app-api/control/file/download/:fileName"])
}

/*
* 测试用例描述：TestFileController_Prepare
* 预置条件：
* 操作步骤：
*     1. 创建 FileController 实例
*     2. 调用 Prepare 方法
* 预期结果：
*     1. fileService 被正确初始化
* 修改历史：
*     1. 2025-8-26 新建测试用例
 */
func TestFileControllerPrepare(t *testing.T) {
	controller := &FileController{}
	controller.Prepare()

	assert.NotNil(t, controller.fileService)
}
