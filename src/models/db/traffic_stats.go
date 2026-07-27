// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package db

import (
	"github.com/beego/beego/v2/client/orm"
)

type MediaTrafficStats struct {
	ID         int    `json:"-" orm:"auto;pk;column(id)"`
	SessionID  string `json:"session_id" orm:"column(session_id)"`
	AppType    int    `json:"app_type" orm:"column(app_type)"`
	StartedAt  string `json:"started_at" orm:"column(started_at)"`
	FinishedAt string `json:"finished_at" orm:"column(finished_at)"`
	OutBytes   int64  `json:"out_bytes" orm:"column(out_bytes)"`
	AccessType int    `json:"access_type" orm:"column(access_type)"`
}

func (s *MediaTrafficStats) TableName() string {
	return "t_media_traffic_stats"
}

func (s *MediaTrafficStats) Validate() error {
	return nil
}

type ControlTrafficStats struct {
	ID         int    `json:"-" orm:"auto;pk;column(id)"`
	SessionID  string `json:"session_id" orm:"column(session_id)"`
	AppType    int    `json:"app_type" orm:"column(app_type)"`
	StartedAt  string `json:"started_at" orm:"column(started_at)"`
	FinishedAt string `json:"finished_at,omitempty" orm:"column(finished_at)"`
	OutBytes   int64  `json:"out_bytes,omitempty" orm:"column(out_bytes)"`
	AccessType int    `json:"access_type" orm:"column(access_type)"`
}

func (s *ControlTrafficStats) TableName() string {
	return "t_control_traffic_stats"
}

func (s *ControlTrafficStats) Validate() error {
	return nil
}

type SessionStats struct {
	ID          int    `json:"-" orm:"auto;pk;column(id)"`
	SessionID   string `json:"session_id" orm:"column(session_id)"`
	AppType     int    `json:"app_type" orm:"column(app_type)"`
	StartedAt   string `json:"started_at" orm:"column(started_at)"`
	FinishedAt  string `json:"finished_at" orm:"column(finished_at)"`
	TcpUniqueId string `json:"tcp_unique_id" orm:"column(tcp_unique_id)"`
}

func (s *SessionStats) TableName() string {
	return "t_session_stats"
}

func (s *SessionStats) Validate() error {
	return nil
}

func init() {
	orm.RegisterModel(&MediaTrafficStats{})
	orm.RegisterModel(&ControlTrafficStats{})
	orm.RegisterModel(&SessionStats{})
}
