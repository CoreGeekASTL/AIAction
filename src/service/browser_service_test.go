/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package service

import (
	"testing"

	"Go-chassis-extend/api/GSF/api/base"
	"github.com/stretchr/testify/assert"

	"GIDS/common/cse"
	"GIDS/models/browsergateway"
)

const (
	testBrowserInnerEndpoint  = "10.0.0.1:8090"
	testMediaInnerEndpoint    = "10.0.0.1:30002"
	testMediaExtendEndpoint   = "10.0.0.2:20010"
	testControlExtendEndpoint = "10.0.0.2:20009"
)

/*
* 测试用例描述：TestBrowserServiceImplGetAllReadyServiceInstances
* 预置条件：
* 操作步骤：
*		1. 注册中心没有上报实例
*	  	2. 注册中心有一个ready实例
*		3. 注册中心有实例，且插件加载失败
*		4. 注册中心有实例，但cap为0
* 预期结果：
*     	1. 返回nil
*		2. 返回一个实例
*		3. 返回nil
*		3. 返回nil
* 修改历史：
*     1. 2025-11-07 新建测试用例
*     2. 2025-11-14 修改测试用例
 */
func TestBrowserServiceImplGetAllReadyServiceInstances(t *testing.T) {

	tests := []struct {
		name string
		cse  cse.Cse
		want []browsergateway.ServiceInstance
	}{
		{
			name: "no instances",
			cse:  &noInstanceCse{},
			want: nil,
		},
		{
			name: "one ready",
			cse:  &oneReadyInstanceCse{},
			want: []browsergateway.ServiceInstance{{BrowserInnerEndpoint: testBrowserInnerEndpoint,
				MediaInnerEndpoint:    testMediaInnerEndpoint,
				MediaExtendEndpoint:   testMediaExtendEndpoint,
				ControlExtendEndpoint: testControlExtendEndpoint,
				Cap:                   40,
				Used:                  1,
				PluginStatus:          "Completed",
				IsHealthy:             true,
			}},
		},
		{
			name: "plugin not completed",
			cse:  &notCompleteCse{},
			want: nil,
		},
		{
			name: "cap equals to 0",
			cse:  &capIllegalCse{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BrowserServiceImpl{
				cse: tt.cse,
			}
			assert.Equalf(t, tt.want, b.GetAllReadyServiceInstances(), "GetAllReadyServiceInstances()")
		})
	}
}

type noInstanceCse struct {
}

func (c *noInstanceCse) GetBrowserGateWayInstanceByInnerEndpoint(innerEndpoint string) (browsergateway.ServiceInstance, bool) {
	return browsergateway.ServiceInstance{}, false
}

func (c *noInstanceCse) GetAllBrowserGateWayInstances() []browsergateway.ServiceInstance {
	return nil
}

func (c *noInstanceCse) GetAllMicroServiceInstanceInfo(serviceName string) ([]base.MicroServiceInstance, error) {
	return nil, nil
}

func (c *noInstanceCse) AddChainEndpoint(endpoint string) {}

func (c *noInstanceCse) Report(maxRetry int) {}

type oneReadyInstanceCse struct {
}

func (c *oneReadyInstanceCse) GetBrowserGateWayInstanceByInnerEndpoint(innerEndpoint string) (browsergateway.ServiceInstance, bool) {
	return browsergateway.ServiceInstance{}, false
}

func (c *oneReadyInstanceCse) GetAllBrowserGateWayInstances() []browsergateway.ServiceInstance {
	return []browsergateway.ServiceInstance{{BrowserInnerEndpoint: testBrowserInnerEndpoint,
		MediaInnerEndpoint:    testMediaInnerEndpoint,
		MediaExtendEndpoint:   testMediaExtendEndpoint,
		ControlExtendEndpoint: testControlExtendEndpoint,
		Cap:                   40,
		Used:                  1,
		PluginStatus:          "Completed",
		IsHealthy:             true,
	}}
}

func (c *oneReadyInstanceCse) GetAllMicroServiceInstanceInfo(serviceName string) ([]base.MicroServiceInstance, error) {
	return nil, nil
}

func (c *oneReadyInstanceCse) AddChainEndpoint(endpoint string) {}

func (c *oneReadyInstanceCse) Report(maxRetry int) {}

type notCompleteCse struct {
}

func (c *notCompleteCse) GetBrowserGateWayInstanceByInnerEndpoint(innerEndpoint string) (browsergateway.ServiceInstance, bool) {
	return browsergateway.ServiceInstance{}, false
}

func (c *notCompleteCse) GetAllBrowserGateWayInstances() []browsergateway.ServiceInstance {
	return []browsergateway.ServiceInstance{{BrowserInnerEndpoint: testBrowserInnerEndpoint,
		MediaInnerEndpoint:    testMediaInnerEndpoint,
		MediaExtendEndpoint:   testMediaExtendEndpoint,
		ControlExtendEndpoint: testControlExtendEndpoint,
		Cap:                   40,
		Used:                  1,
		PluginStatus:          "Failed",
	}}
}

func (c *notCompleteCse) GetAllMicroServiceInstanceInfo(serviceName string) ([]base.MicroServiceInstance, error) {
	return nil, nil
}

func (c *notCompleteCse) AddChainEndpoint(endpoint string) {}

func (c *notCompleteCse) Report(maxRetry int) {}

type capIllegalCse struct {
}

func (c *capIllegalCse) GetBrowserGateWayInstanceByInnerEndpoint(innerEndpoint string) (browsergateway.ServiceInstance, bool) {
	return browsergateway.ServiceInstance{}, false
}

func (c *capIllegalCse) GetAllBrowserGateWayInstances() []browsergateway.ServiceInstance {
	return []browsergateway.ServiceInstance{{BrowserInnerEndpoint: testBrowserInnerEndpoint,
		MediaInnerEndpoint:    testMediaInnerEndpoint,
		MediaExtendEndpoint:   testMediaExtendEndpoint,
		ControlExtendEndpoint: testControlExtendEndpoint,
		Cap:                   40,
		Used:                  1,
		PluginStatus:          "Failed",
	}}
}

func (c *capIllegalCse) GetAllMicroServiceInstanceInfo(serviceName string) ([]base.MicroServiceInstance, error) {
	return nil, nil
}

func (c *capIllegalCse) AddChainEndpoint(endpoint string) {}

func (c *capIllegalCse) Report(maxRetry int) {}
