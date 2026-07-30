/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.huawei.com/fusionstage/greatwall-sdk-go/overloadcontroller"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
)

const testRateLimit = 10

var policy = `{
  "APIService": {
    "default": [
      {
        "type": "concurrent_limit",
        "mode": "local",
        "time_unit": "second",
        "time_interval": 1,
        "rate_limit": %d
      }
    ]
  }
}`

/*
* 测试用例描述： TestOverLoadFilterGranted
* 预置条件：无
* 操作步骤：
*     1. 初始化限流策略
*     2. 调用OverLoadFilter
*     3. 检测context状态码
* 预期结果：
*     1. 初始化策略成功
*     2. context状态码不为429
* 修改历史：
*     1. 2025-12-10 新建测试用例
 */
func TestOverLoadFilterGranted(t *testing.T) {
	err := overloadcontroller.InitWithPolicies(fmt.Sprintf(policy, testRateLimit))
	assert.NoError(t, err)
	context := beecontext.NewContext()
	context.Input = beecontext.NewInput()
	context.Input.Context = context
	context.Request, err = http.NewRequest(http.MethodGet, "/api", nil)
	assert.NoError(t, err)
	context.ResponseWriter = &beecontext.Response{}
	OverLoadFilter(context)
	assert.NotEqual(t, http.StatusTooManyRequests, context.Output.Status)
}

/*
* 测试用例描述： TestOverLoadFilterNotGranted
* 预置条件：无
* 操作步骤：
*     1. 初始化限流策略
*     2. 调用OverLoadFilter
*     3. 检测context状态码
* 预期结果：
*     1. 初始化策略成功
*     2. context状态码为429
* 修改历史：
*     1. 2025-12-10 新建测试用例
 */
func TestOverLoadFilterNotGranted(t *testing.T) {
	err := overloadcontroller.InitWithPolicies(fmt.Sprintf(policy, 0))
	assert.NoError(t, err)
	context := beecontext.NewContext()
	context.Input = beecontext.NewInput()
	context.Input.Context = context
	context.Request, err = http.NewRequest(http.MethodGet, "/api", nil)
	assert.NoError(t, err)
	context.ResponseWriter = &beecontext.Response{ResponseWriter: httptest.NewRecorder()}
	OverLoadFilter(context)
	assert.Equal(t, http.StatusTooManyRequests, context.ResponseWriter.Status)
}
