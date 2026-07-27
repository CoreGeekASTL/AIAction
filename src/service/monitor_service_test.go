/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package service

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"GIDS/models/monitor"
)

/*
* 测试用例描述： TestMonitorSchedule
* 预置条件：无
* 操作步骤：
*     1. 测试能正常获取指标的值
*     2. 测试获取指标没有值
*     3. 测试没有用于获取指标值的方法
* 预期结果：
*     1. 返回指标对象 mocIdMap非空
*     2. 返回指标对象 mocIdMap为空
*     3. 返回指标对象 mocIdMap为空
* 修改历史：
*     1. 2025-12-14 新建测试用例
 */
const testMocIDValue = 1003

var testMonitorMocID = monitor.MocID(testMocIDValue)

func TestMonitorSchedule(t *testing.T) {
	tests := []struct {
		name           string
		monitorService *MonitorServiceImpl
		monitorConfig  *monitor.MonitorConfig
		want           struct {
			mocMapLen int
		}
	}{
		{
			name: "MetricOnlineUsersNormal",
			monitorService: &MonitorServiceImpl{
				monitorConfig: monitor.MonitorConfig{
					MetricGroups: []monitor.MetricGroup{
						{
							GroupID: "3001",
							MocId:   testMonitorMocID,
							Metrics: []monitor.Metric{
								{
									ID: monitor.MetricOnlineUsers, // normal query
								},
							},
						},
					},
				},
				mocIdMap: make(map[monitor.MocID]map[string]struct{}, 5),
				statsService: &TrafficStatsServiceMock{
					onlineRes: []Res{
						{Obj: "default_obj", Cnt: 66},
						{Obj: "", Cnt: 66.1},
					},
				},
			},
			want: struct{ mocMapLen int }{
				mocMapLen: 1,
			},
		},
		{
			name: "MetricOnlineUsersNoRes",
			monitorService: &MonitorServiceImpl{
				monitorConfig: monitor.MonitorConfig{
					MetricGroups: []monitor.MetricGroup{
						{
							GroupID: "3001",
							MocId:   testMonitorMocID,
							Metrics: []monitor.Metric{
								{
									ID: monitor.MetricOnlineUsersPerModel, // normal query
								},
							},
						},
					},
				},
				mocIdMap: make(map[monitor.MocID]map[string]struct{}, 5),
				statsService: &TrafficStatsServiceMock{
					onlineOfModelRes: []Res{},
				},
			},
			want: struct{ mocMapLen int }{
				mocMapLen: 0,
			},
		},
		{
			name: "MetricOnlineUsersNoFunc",
			monitorService: &MonitorServiceImpl{
				monitorConfig: monitor.MonitorConfig{
					MetricGroups: []monitor.MetricGroup{
						{
							GroupID: "3001",
							MocId:   testMonitorMocID,
							Metrics: []monitor.Metric{
								{
									ID: 12345, // normal query
								},
							},
						},
					},
				},
				mocIdMap: make(map[monitor.MocID]map[string]struct{}, 5),
				statsService: &TrafficStatsServiceMock{
					trafficOfAppRes: []Res{
						{Obj: "obj_001", Cnt: 66},
					},
				},
			},
			want: struct{ mocMapLen int }{
				mocMapLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.monitorService.createMetricFunctionMap()
			tt.monitorService.monitorSchedule()
			mocIdMapLen := len(tt.monitorService.mocIdMap)
			assert.Equal(t, mocIdMapLen, tt.want.mocMapLen)
		})
	}
}
