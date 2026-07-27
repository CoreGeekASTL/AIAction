/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package cse 提供CSE微服务注册、发现与实例管理功能。
package cse

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"GIDS/common/logger"
	"GIDS/models/browsergateway"
	"Go-chassis-extend/api/GSF/api"
	"Go-chassis-extend/api/GSF/api/base"
	"Go-chassis-extend/api/ServiceComb/go-chassis/core/config"
)

type Cse interface {
	GetAllMicroServiceInstanceInfo(serviceName string) ([]base.MicroServiceInstance, error)
	GetAllBrowserGateWayInstances() []browsergateway.ServiceInstance
	GetBrowserGateWayInstanceByInnerEndpoint(innerEndpoint string) (browsergateway.ServiceInstance, bool)
	AddChainEndpoint(endpoint string)
	Report(maxRetry int)
}

type cse struct {
	register           api.Registry
	appid              string
	browserGWInstances sync.Map
	chainEndpoints     []string
}

// GetBrowserGateWayInstanceByInnerEndpoint get browserGateway instance by inner endpoint
func (c *cse) GetBrowserGateWayInstanceByInnerEndpoint(innerEndpoint string) (browsergateway.ServiceInstance, bool) {
	targetInstance := browsergateway.ServiceInstance{}
	c.browserGWInstances.Range(func(k, v any) bool {
		instance, ok := v.(browsergateway.ServiceInstance)
		if !ok {
			return true
		}
		if instance.BrowserInnerEndpoint == innerEndpoint {
			targetInstance = instance
			return false
		}
		return true
	})
	return targetInstance, targetInstance.BrowserInnerEndpoint == innerEndpoint
}

// GetAllMicroServiceInstanceInfo get all service instances for register
func (c *cse) GetAllMicroServiceInstanceInfo(serviceName string) ([]base.MicroServiceInstance, error) {
	msKey := base.MicroServiceKey{
		AppId:       c.appid,
		ServiceName: serviceName,
		Version:     "0+",
	}
	selfServiceID := config.GetSelfServiceID()
	return c.register.GetAllMicroServiceInstanceInfo(selfServiceID, msKey)
}

var cseService cse

func NewCse() Cse {
	return &cseService
}

func Init() {
	cseService = cse{appid: os.Getenv("APPID"),
		register:           api.NewRegistry(),
		browserGWInstances: sync.Map{},
		chainEndpoints:     make([]string, 0, 2),
	}
	selfServiceID := config.GetSelfServiceID()
	msKey := base.MicroServiceKey{
		AppId:       cseService.appid,
		ServiceName: "browser-gateway",
		Version:     "0+",
	}
	err := cseService.register.WatchMicroServiceV1(selfServiceID, []base.MicroServiceKey{msKey}, browserGWNotifier{})
	if err != nil {
		logger.Errorf("watch failed : %v", err)
	} else {
		logger.Infof("Watch success")
	}
}

type browserGWNotifier struct {
}

// WatchServiceCallBack callback func
func (c browserGWNotifier) WatchServiceCallBack(event *base.MicroServiceInstanceChangedEvent) {
	cseService.watchServiceCallBack(event)
}

func (c *cse) watchServiceCallBack(event *base.MicroServiceInstanceChangedEvent) {
	logger.Infof("[watchServiceCallBack] event : %v, browserGWInstances : %v, browserGWInstances list: %v", event, event.Instance, event.InstanceList)
	switch event.Action {
	case "CREATE":
		c.updateInstance(event.Instance)
	case "UPDATE":
		c.updateInstance(event.Instance)
	case "DELETE":
		c.browserGWInstances.Delete(event.Instance.InstanceId)
	case "LIST":
		for i := range event.InstanceList {
			c.updateInstance(event.InstanceList[i].Instance)
		}
	default:
		logger.Infof("unknown event action: %s", event.Action)
	}
}

func (c *cse) updateInstance(m *base.MicroServiceInstance) {
	statusStr := m.Properties["status"]
	if statusStr == "" {
		return
	}
	instance := browsergateway.ServiceInstance{}
	err := json.Unmarshal([]byte(statusStr), &instance)
	if err != nil {
		logger.Errorf("[WatchServiceCallBack] failed to unmarshal %s: %v", statusStr, err)
		return
	}
	isHealthy, ok := m.Properties["isHealthy"]
	if !ok {
		instance.IsHealthy = true
	} else {
		isHealthyBool, err := strconv.ParseBool(isHealthy)
		if err != nil {
			logger.Errorf("[updateInstance] parse isHealthy failed: %v", err)
			instance.IsHealthy = true
		} else {
			instance.IsHealthy = isHealthyBool
		}
	}

	instance.CheckMsg = m.Properties["checkMsg"]

	c.browserGWInstances.Store(m.InstanceId, instance)
}

// GetAllBrowserGateWayInstances get all browserGateway instances
func (c *cse) GetAllBrowserGateWayInstances() []browsergateway.ServiceInstance {
	instances := make([]browsergateway.ServiceInstance, 0, 10)
	c.browserGWInstances.Range(func(key, value any) bool {
		instance, ok := value.(browsergateway.ServiceInstance)
		if ok {
			instances = append(instances, instance)
		}
		return true
	})
	return instances
}
func (c *cse) AddChainEndpoint(endpoint string) {
	c.chainEndpoints = append(c.chainEndpoints, endpoint)
}

const (
	retryIntervalSec = 30
	retryStepLen     = 1
)

// Report will upload chain endpoints to cs
// Report will upload chain endpoints to cse
func (c *cse) Report(maxRetry int) {
	selfServiceID := config.GetSelfServiceID()
	selfInstanceID := config.GetSelfInstanceID()

	m := map[string]string{
		"chainEndpoints": strings.Join(c.chainEndpoints, ","),
	}
	if maxRetry <= 0 {
		logger.Fatalf("exceed max retry to report chain endpoints, exit")
	}
	err := c.register.UpdateMicroServiceInstanceProperties(selfServiceID, selfInstanceID, m)
	if err != nil {
		logger.Errorf("failed to update user count propertiy, err: %v, will retry", err)
		time.Sleep(retryIntervalSec * time.Second)
		c.Report(maxRetry - retryStepLen)
	}
}
