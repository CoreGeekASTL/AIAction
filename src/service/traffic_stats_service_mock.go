/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package service

import (
	"encoding/json"

	"GIDS/models/db"
)

var _ TrafficStatsService = &TrafficStatsServiceMock{}

// TrafficStatsServiceMock 实现 TrafficStatsService 接口
type TrafficStatsServiceMock struct {
	onlineRes             []Res
	onlineOfModelRes      []Res
	trafficOfAppRes       []Res
	onlineOfInstanceRes   []Res
	trafficRes            []Res
	batchInsertRes        []Res
	HandleSessionStatsRes []Res
}

// GetOnline 获取在线用户数
func (s *TrafficStatsServiceMock) GetOnline(_, _ string) []Res {
	return s.onlineRes
}

// GetOnlineOfModel 获取各机型的在线用户数
func (s *TrafficStatsServiceMock) GetOnlineOfModel(_, _ string) []Res {
	return s.onlineOfModelRes
}

// GetOnlineOfInstance 获取各个实例的在线用户数
func (s *TrafficStatsServiceMock) GetOnlineOfInstance(_, _ string) []Res {
	return s.onlineOfInstanceRes
}

// GetTrafficOfApp 获取个应用流量
func (s *TrafficStatsServiceMock) GetTrafficOfApp(_, _ string) []Res {
	return s.trafficOfAppRes
}

// GetTraffic 获取站点/总流量
func (s *TrafficStatsServiceMock) GetTraffic(_, _ string) []Res {
	return s.trafficRes
}

// ExportSessionList 导出会话数据到 CSV 文件
func (s *TrafficStatsServiceMock) ExportSessionList(month string, csvFilePath string) error {
	return nil
}

// ExportMediaList 导出媒体流数据到 CSV 文件
func (s *TrafficStatsServiceMock) ExportMediaList(month string, csvFilePath string) error {
	return nil
}

// ExportControlList 导出控制流数据到 CSV 文件
func (s *TrafficStatsServiceMock) ExportControlList(month string, csvFilePath string) error {
	return nil
}

// 实现BatchInsertStats 方法
func (s *TrafficStatsServiceMock) BatchInsertStats(tag string, items []json.RawMessage) error {
	return nil
}

func (s *TrafficStatsServiceMock) HandleSessionStats(session *db.SessionStats) error { return nil }
