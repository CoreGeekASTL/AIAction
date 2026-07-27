/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package cse

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"Go-chassis-extend/api/GSF/api/base"
	"Go-chassis-extend/api/ServiceComb/go-chassis/core/bundle/dependency"
	registry2 "Go-chassis-extend/api/ServiceComb/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"

	"GIDS/models/browsergateway"
)

/*
* 测试用例描述：测试 GetAllMicroServiceInstanceInfo 函数
* 预置条件：
* 操作步骤：
*     1. 使用MockRegistry调用 GetAllMicroServiceInstanceInfo函数
* 预期结果：
*     1. 长度为1，InstanceId为instance1
* 修改历史：
*     1. 2025-12-24 新建测试用例
 */
func TestGetAllMicroServiceInstanceInfo(t *testing.T) {
	cseService.register = &mockRegistry{}

	instances, err := cseService.GetAllMicroServiceInstanceInfo("browser-gateway")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(instances))
	assert.Equal(t, "instance1", instances[0].InstanceId)
}

/*
* 测试用例描述：测试 GetBrowserGateWayInstanceByInnerEndpoint 函数
* 预置条件：browserGWInstances 中包含一个实例，其 BrowserInnerEndpoint 为 "endpoint1"
* 操作步骤：
*     1. 调用 GetBrowserGateWayInstanceByInnerEndpoint 函数，传入 "endpoint1"
* 预期结果：
*     1. 找到实例，返回的实例 BrowserInnerEndpoint 为 "endpoint1"
* 修改历史：
*     1. 2025-12-24 新建测试用例
*	  2. 2026-01-04 修改测试用例
 */
func TestGetBrowserGateWayInstanceByInnerEndpoint(t *testing.T) {
	cseService.browserGWInstances = sync.Map{}
	cseService.browserGWInstances.Store("instance1", browsergateway.ServiceInstance{BrowserInnerEndpoint: "endpoint1"})

	instance, found := cseService.GetBrowserGateWayInstanceByInnerEndpoint("endpoint1")
	assert.True(t, found)
	assert.Equal(t, "endpoint1", instance.BrowserInnerEndpoint)
}

/*
* 测试用例描述：测试 GetAllBrowserGateWayInstances 函数
* 预置条件：browserGWInstances 中包含两个实例，分别具有不同的 BrowserInnerEndpoint
* 操作步骤：
*     1. 调用 GetAllBrowserGateWayInstances 函数
* 预期结果：
*     1. 返回的实例列表长度为 2
* 修改历史：
*     1. 2025-12-24 新建测试用例
*	  2. 2026-01-04 修改测试用例
 */
const expectedInstanceCount = 2

func TestGetAllBrowserGateWayInstances(t *testing.T) {
	cseService.browserGWInstances = sync.Map{}
	cseService.browserGWInstances.Store("instance1", browsergateway.ServiceInstance{BrowserInnerEndpoint: "endpoint1"})
	cseService.browserGWInstances.Store("instance2", browsergateway.ServiceInstance{BrowserInnerEndpoint: "endpoint2"})

	instances := cseService.GetAllBrowserGateWayInstances()
	assert.Equal(t, expectedInstanceCount, len(instances))
}

/*
* 测试用例描述：测试 watchServiceCallBack 函数
* 预置条件：browserGWInstances 中包含两个实例
* 操作步骤：
*     1. 调用 watchServiceCallBack 函数，传入不同类型的事件（CREATE 和 DELETE）
* 预期结果：
*     1. CREATE 事件后，browserGWInstances 的长度增加
*     2. DELETE 事件后，browserGWInstances 的长度减少
* 修改历史：
*     1. 2025-12-24 新建测试用例
*	  2. 2026-01-04 修改测试用例
 */
func TestWatchServiceCallBack(t *testing.T) {
	cseService.browserGWInstances = sync.Map{}
	cseService.browserGWInstances.Store("instance1", browsergateway.ServiceInstance{BrowserInnerEndpoint: "endpoint1"})
	cseService.browserGWInstances.Store("instance2", browsergateway.ServiceInstance{BrowserInnerEndpoint: "endpoint2"})

	tests := []struct {
		name    string
		event   *base.MicroServiceInstanceChangedEvent
		wantLen int
	}{
		{
			name: "create event",
			event: &base.MicroServiceInstanceChangedEvent{
				Action: "CREATE",
				Instance: &base.MicroServiceInstance{
					InstanceId: "instance3",
					Properties: map[string]string{
						"status":    `{"browserGWInnerEndpoint": "endpoint3"}`,
						"isHealthy": "true",
						"checkMsg":  "healthy",
					},
				},
			},
			wantLen: 3,
		},
		{
			name: "DELETE event",
			event: &base.MicroServiceInstanceChangedEvent{
				Action: "DELETE",
				Instance: &base.MicroServiceInstance{
					InstanceId: "instance3",
				},
			},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cseService.watchServiceCallBack(tt.event)
			instancesLen := 0
			cseService.browserGWInstances.Range(func(key, value any) bool {
				instancesLen++
				return true
			})
			assert.Equal(t, tt.wantLen, instancesLen)
		})
	}
}

/*
* 测试用例描述：测试 updateInstance 函数
* 预置条件：browserGWInstances 为空
* 操作步骤：
*     1. 调用 updateInstance 函数，传入一个包含状态信息的实例
* 预期结果：
*     1. browserGWInstances 中包含该实例，且实例信息正确
* 修改历史：
*     1. 2025-12-24 新建测试用例
*	  2. 2026-01-04 修改测试用例
 */
func TestUpdateInstance(t *testing.T) {
	cseService.browserGWInstances = sync.Map{}
	instance := &base.MicroServiceInstance{
		InstanceId: "instance4",
		Properties: map[string]string{
			"status":    `{"browserGWInnerEndpoint": "endpoint4"}`,
			"isHealthy": "false",
			"checkMsg":  "unhealthy",
		},
	}

	cseService.updateInstance(instance)
	instancesLen := 0
	cseService.browserGWInstances.Range(func(key, value any) bool {
		instancesLen++
		return true
	})
	assert.Equal(t, 1, instancesLen)
	value, ok := cseService.browserGWInstances.Load("instance4")
	assert.Equal(t, true, ok)
	assert.Equal(t, "endpoint4", value.(browsergateway.ServiceInstance).BrowserInnerEndpoint)
}

type mockRegistry struct{}

func (m *mockRegistry) GetAllMicroServiceInstanceInfo(serviceID string, msKey base.MicroServiceKey) ([]base.MicroServiceInstance, error) {
	return []base.MicroServiceInstance{
		{
			InstanceId: "instance1",
			Properties: map[string]string{
				"status":    `{"BrowserInnerEndpoint": "endpoint1"}`,
				"isHealthy": "true",
				"checkMsg":  "healthy",
			},
		},
	}, nil
}

func (m *mockRegistry) WatchMicroServiceV1(selfMicroServiceId string, providers []base.MicroServiceKey, callback base.WatchServiceNotify, opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetAllMicroServices() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) RegisterServiceAndInstanceInfo(microService *base.MicroService, instance *base.MicroServiceInstance) (string, string, error) {
	return "", "", fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetMicroServiceId(appId, microServiceName, version string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetMicroServiceIds(appId, microServiceName string) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetMicroServiceInstances(consumerId, providerId string) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetMicroServiceInstanceInfo(consumerId, providerId string) ([]base.MicroServiceInstance, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) FindMicroServiceInstances(consumerId string, microService base.MicroServiceKey) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) UnregisterMicroServiceInstance(microServiceId, microServiceInstanceId string) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) WatchMicroService(selfMicroServiceId string, providers []base.MicroServiceKey, callback func(*base.MicroServiceInstanceChangedEvent), opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) UpdateMicroServiceInstanceStatus(microServiceId, microServiceInstanceId, status string) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) UpdateMicroServiceInstanceProperties(microServiceId, microServiceInstanceId string, Properties map[string]string) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) WatchMicroServiceProperties(selfServiceId string, properties map[string]string, callback func(map[string]string, *base.MicroServiceInstanceChangedEvent)) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) WatchMicroServiceInstnaceProperties(selfServiceId string, properties map[string]string, callback func(map[string]string, *base.MicroServiceInstanceChangedEvent), opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) WatchMicroServiceInstancePropertiesV1(selfServiceId string, properties map[string]string, callback base.WatchLabelNotify, opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetAllMicroServiceInstancesByProperty(consumerId string, propertis map[string]string) ([]*base.MicroServiceInstance, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetAllMicroServiceInstancesByPropertyAppID(consumerId, appID string, propertis map[string]string) ([]*base.MicroServiceInstance, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetInstanceByAppId(appId string) (map[base.MicroServiceKey][]base.MicroServiceInstance, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) UnWatchInstanceProperties(selfServiceId string, properties map[string]string, callback func(map[string]string, *base.MicroServiceInstanceChangedEvent), opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) UnWatchMicroServicePropertiesV1(selfServiceId string, properties map[string]string, callback base.WatchLabelNotify, opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) UnWatchMicroService(selfMicroServiceId string, providers []base.MicroServiceKey, callback func(*base.MicroServiceInstanceChangedEvent), opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) UnWatchMicroServiceV1(selfMicroServiceId string, providers []base.MicroServiceKey, callback base.WatchServiceNotify, opt ...dependency.WatchOption) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) AppendDependency(cDep *registry2.MicroServiceDependency) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) RegisterService(microservice *base.MicroService) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (m *mockRegistry) RegisterServiceInstance(sid string, instance *base.MicroServiceInstance) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (m *mockRegistry) FindMicroServiceInstancesV2(consumerID, appID, microServiceName, version string) ([]*base.MicroServiceInstance, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetAllMicroServicesV2() ([]*base.MicroService, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetMicroService(microServiceID string) (*base.MicroService, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) GetMicroServiceInstancesV2(consumerID, providerID string) ([]*base.MicroServiceInstance, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockRegistry) Close() error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) String() string {
	return "mockRegistry"
}

func (m *mockRegistry) AutoSync() {
}

func (m *mockRegistry) Heartbeat(ctx context.Context, microServiceID, microServiceInstanceID string) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (m *mockRegistry) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) DeleteMicroServiceInstance(microServiceId, microServiceInstanceId string) error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) CSPResClear() error {
	return fmt.Errorf("not implemented")
}

func (m *mockRegistry) ForceListEvent() error {
	return fmt.Errorf("not implemented")
}
