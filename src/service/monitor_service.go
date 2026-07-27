// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package service
package service

import (
	monitorsdk "CSPGoMonitorSDK/api/monitor"
	gsfConfig "Go-chassis-extend/api/ServiceComb/go-chassis/core/config"
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"GIDS/common/logger"
	"GIDS/models/monitor"
	"GIDS/utils/monitorutil"
)

// MonitorService monitor
type MonitorService interface {
	// InitMonitorSchedule 启动定时任务： CSPGoMonitorSDK 注册完调用该方法； 默认5min定时(跟监控模板填写时保持一致)
	InitMonitorSchedule()
}

func NewMonitorService() MonitorService {
	return &MonitorServiceImpl{}
}

const (
	InitMonitorPeriod = 10 * time.Second
	DotPeriodFiveMin  = 5 * time.Minute

	defaultMonitorFile = "/opt/csp/gids/module/conf/monitor.json"
	defaultSqlYamlFile = "/opt/csp/gids/module/conf/sql.yaml"
)

type getMetricFunc func(startTime, endTime string) []Res

// MonitorServiceImpl monitor report to csp
type MonitorServiceImpl struct {
	metricMapLock sync.RWMutex
	mocIdMap      map[monitor.MocID]map[string]struct{}
	monitorConfig monitor.MonitorConfig
	monitorJson   string
	statsService  TrafficStatsService
	metricFunMap  map[monitor.MetricID]getMetricFunc
	stopChan      chan struct{}
}

// InitMonitorSchedule 启动定时任务： 调用 CSPGoMonitorSDK 注册，注册完打点上报数据； 默认5min定时(跟监控模板填写时保持一致)
func (m *MonitorServiceImpl) InitMonitorSchedule() {
	logger.Infof("Begin to init csp monitor.")

	if err := m.buildMonitorConfig(); err != nil {
		logger.Errorf("build monitor config failed : %v", err)
		return
	}

	// 1. 注册csp话统：先注册才能上报数据。每间隔一段时间进行注册，直到注册成功然后停止
	for {
		if err := m.InitCspMonitor(); err != nil {
			logger.Errorf("init csp monitor failed, will try again after 10 seconds. err : %v", err)
			time.Sleep(InitMonitorPeriod)
		} else {
			// 2. 开始csp打点，即定时上报各项指标的数据
			m.startCspMonitor()
			break
		}
	}
}

func (m *MonitorServiceImpl) buildMonitorConfig() error {
	monitorJsonFile := beego.AppConfig.DefaultString("cspmonitor::monitorJsonFile", defaultMonitorFile)
	monitorData, err := os.ReadFile(monitorJsonFile)
	if err != nil {
		return fmt.Errorf("read file %s error, %v", monitorJsonFile, err)
	}

	if err = json.Unmarshal(monitorData, &m.monitorConfig); err != nil {
		return fmt.Errorf("parse json %s error, %v", monitorData, err)
	}
	m.monitorJson = string(monitorData)
	m.mocIdMap = make(map[monitor.MocID]map[string]struct{}, len(m.monitorConfig.MetricGroups))
	return nil
}

func (m *MonitorServiceImpl) InitCspMonitor() error {
	appId := gsfConfig.GetGlobalDefinition().AppID
	serviceName := gsfConfig.GetSelfServiceName()
	instanceName := gsfConfig.GetSelfInstanceID()
	logger.Infof("Init csp monitor sdk, appID : %s, serviceName : %s, instanceName : %s ",
		appId, serviceName, instanceName)

	id, err := strconv.Atoi(appId)
	if err != nil {
		return err
	}

	// 1. 初始化监控sdk
	err = monitorsdk.MonSdkInstance.InitMonitor(id, serviceName, instanceName)
	if err != nil {
		return fmt.Errorf("sdk initMonitor error, %v", err)
	}

	// 2. 模型注册
	err = monitorsdk.MonSdkInstance.RegisterBasicInfo(m.monitorJson)
	if err != nil {
		return fmt.Errorf("sdk registerBasicInfo error, %v", err)
	}
	return nil
}

func (m *MonitorServiceImpl) startCspMonitor() {
	logger.Infof("Start csp monitor")
	sqlYamlFile := beego.AppConfig.DefaultString("cspmonitor::sqlYamlFile", defaultSqlYamlFile)
	m.statsService = NewTrafficStatsService(sqlYamlFile)
	m.createMetricFunctionMap()
	m.stopChan = make(chan struct{})
	go func() {
		ticker := time.NewTicker(DotPeriodFiveMin)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.monitorSchedule()
			case <-m.stopChan:
				logger.Infof("stopping csp monitor")
				return
			}
		}
	}()
}

func (m *MonitorServiceImpl) Stop() {
	if m.stopChan != nil {
		close(m.stopChan)
	}
}

// 创建指标函数映射
func (m *MonitorServiceImpl) createMetricFunctionMap() {
	m.metricFunMap = map[monitor.MetricID]getMetricFunc{
		monitor.MetricOnlineUsers:         m.statsService.GetOnline,
		monitor.MetricOnlineUsersPerModel: m.statsService.GetOnlineOfModel,
		monitor.MetricUsersSupportedByVm:  m.statsService.GetOnlineOfInstance,
		monitor.MetricApplicationTraffic:  m.statsService.GetTrafficOfApp,
		monitor.MetricSiteTraffic:         m.statsService.GetTraffic,
	}
}

func (m *MonitorServiceImpl) monitorSchedule() {
	logger.Infof("Start csp monitor schedule")
	for _, group := range m.monitorConfig.MetricGroups {
		for _, metric := range group.Metrics {
			logger.Infof("Start schedule mocid is %d, moiid is %d", group.MocId, metric.ID)

			metricFunc, exists := m.metricFunMap[metric.ID]
			if !exists || metricFunc == nil {
				logger.Errorf(" err get metric of metric %d : metricFunc not exist ", metric.ID)
				continue
			}

			// 1、 获取指标的结果 ， 例如各机型在线人数(指标): [{机型AAA, 34}, {机型BBB, 34}, ... ]
			results := m.getMetricResults(metricFunc, metric.ID)
			if results == nil || len(results) == 0 {
				logger.Warnf("No data for metric %d", metric.ID)
				continue
			}

			//  2、 处理指标的结果
			m.processMetricResults(group.MocId, metric, results)
		}
	}
}

// 执行指标函数并捕获panic
func (m *MonitorServiceImpl) getMetricResults(metricFunc getMetricFunc, metricID monitor.MetricID) []Res {
	var results []Res
	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("panic occurred in metricFunc for metric %d: %v", metricID, r)
				logger.Errorf("stack trace: %s", debug.Stack())
			}
		}()
		results = metricFunc(monitorutil.GetLastFiveMinuteWindow(nil))
	}()
	return results
}

// 处理单个指标项
func (m *MonitorServiceImpl) processMetricResults(mocId monitor.MocID, metric monitor.Metric, results []Res) {
	for _, item := range results {
		if item.Obj == "" {
			logger.Errorf("empty obj for metric %d", metric.ID)
			continue
		}
		if err := m.addMoiIdIfNotExists(mocId, item.Obj); err != nil {
			logger.Errorf("err add moiId of metric %d: %v", metric.ID, err)
			continue
		}
		if err := monitorsdk.MonSdkInstance.SetMetric(int(metric.ID), item.Obj, float64(item.Cnt)); err != nil {
			logger.Errorf("err set metric of metric %d: %v", metric.ID, err)
			continue
		}
	}
}

func (m *MonitorServiceImpl) addMoiIdIfNotExists(mocId monitor.MocID, moiId string) error {
	m.metricMapLock.Lock()
	defer m.metricMapLock.Unlock()
	if _, exists := m.mocIdMap[mocId]; !exists {
		m.mocIdMap[mocId] = make(map[string]struct{})
	}

	if _, exists := m.mocIdMap[mocId][moiId]; !exists {
		if err := monitorsdk.MonSdkInstance.ObjChange(uint32(mocId), 1, moiId); err != nil {
			return fmt.Errorf("err register moi %s for moc %d", moiId, mocId)
		} else {
			m.mocIdMap[mocId][moiId] = struct{}{}
		}
	}
	return nil
}
