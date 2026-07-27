// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package monitor
package monitor

import (
	"encoding/json"
	"strconv"
)

type Metric struct {
	ID           MetricID `json:"ID"`
	ValueType    string   `json:"ValueType"`
	MeType       string   `json:"MeType"`
	DefaultValue string   `json:"DefaultValue"`
}

type MetricGroup struct {
	GroupID    string   `json:"GroupID"`
	MocId      MocID    `json:"MocId"`
	MocType    string   `json:"MocType"`
	MoiMgrType string   `json:"MoiMgrType"`
	Metrics    []Metric `json:"Metrics"`
}

// MonitorConfig 结构跟 monitor.json 的结构保持一致，然后可减少部分非必要参数
type MonitorConfig struct {
	MetricGroups     []MetricGroup `json:"MetricGroups"`
	RealServicesName string        `json:"RealServicesName"`
}

// UnmarshalJSON 实现 UnmarshalJSON 方法
func (m *MocID) UnmarshalJSON(b []byte) error {
	// 解析为字符串
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	// 将字符串转换为 MocID
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*m = MocID(i)
	return nil
}

func (m *MetricID) UnmarshalJSON(b []byte) error {
	// 解析为字符串
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	// 将字符串转换为 MocID
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*m = MetricID(i)
	return nil
}

type MetricID int
type MocID int

const (
	MetricOnlineUsers         MetricID = 320101
	MetricOnlineUsersPerModel MetricID = 320102
	MetricUsersSupportedByVm  MetricID = 320103
	MetricApplicationTraffic  MetricID = 320201
	MetricSiteTraffic         MetricID = 320202
)
