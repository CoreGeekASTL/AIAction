// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	UploadPath     = "upload"
	DeletePath     = "delete"
	GetAllPath     = "getAll"
	LoadPath       = "load"
	CurrentPath    = "current"
	UploadHandler  = "POST:UploadPluginPackage"
	DeleteHandler  = "POST:DeletePluginPackage"
	GetAllHandler  = "POST:GetPluginPackages"
	LoadHandler    = "POST:LoadPlugin"
	CurrentHandler = "POST:GetCurrentPlugins"
)

/*
* 测试用例描述： TestPluginControllerRouteInfo
* 预置条件：无
* 操作步骤：
*     1. 创建 PluginController 实例
*     2. 调用 RouteInfo 方法
*     3. 检查 RouteInfo 返回的结构体是否为空
* 预期结果：
*     1. 返回正确的 RouteInfo 结构
* 修改历史：
*     1. 2025-8-28 新建测试用例
 */
func TestPluginControllerRouteInfo(t *testing.T) {
	controller := &PluginController{}
	routeInfo := controller.RouteInfo()

	// 检查 routeInfo 是否为空，避免空指针解引用
	if reflect.ValueOf(routeInfo).IsZero() {
		t.Fatalf("RouteInfo is nil")
	}

	ns := "/plugin/v1/"

	// 定义路由映射测试用例
	routeMappingTests := []struct {
		path     string
		expected string
	}{
		{ns + UploadPath, UploadHandler},
		{ns + DeletePath, DeleteHandler},
		{ns + GetAllPath, GetAllHandler},
		{ns + LoadPath, LoadHandler},
		{ns + CurrentPath, CurrentHandler},
	}

	// 遍历测试用例进行断言
	for _, testCase := range routeMappingTests {
		t.Run(testCase.path, func(t *testing.T) {
			assert.Equal(t, testCase.expected, routeInfo.RouteMapping[testCase.path])
		})
	}
}
