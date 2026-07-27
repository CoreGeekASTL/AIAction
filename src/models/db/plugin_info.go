// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package db
package db

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

type PluginPackage struct {
	Field         string       `json:"-" orm:"pk;column(key)"`
	Name          string       `json:"name" orm:"column(name)"`
	Version       string       `json:"version" orm:"column(version)"`
	PackageName   string       `json:"packageName" orm:"column(package_name)"`
	Type          string       `json:"type" orm:"column(plugin_type)"`
	PackageBucket string       `json:"bucket" orm:"column(bucket)"`
	Status        ActiveStatus `json:"status" orm:"column(active_status)"`
	IfActive      bool         `json:"if_active" orm:"column(if_active)"`
	Progress      int          `json:"progress" orm:"column(progress)"`
	CreatedAt     string       `json:"created_at" orm:"column(created_at)"`
}

func (p *PluginPackage) TableName() string {
	return "t_plugin_package"
}

func (p *PluginPackage) GetKey() string {
	return p.GetField()
}

func (p *PluginPackage) GetField() string {
	if p.Field != "" {
		return p.Field
	}
	return fmt.Sprintf("%s:%s:%s", p.Type, p.Name, p.Version)
}

type ActiveStatus string

const (
	Complete ActiveStatus = "Completed"
	Failed   ActiveStatus = "Failed"
	Doing    ActiveStatus = "Doing"
	NotStart ActiveStatus = "NotStart"
)

type PluginActive struct {
	Name     string       `json:"name"`
	Version  string       `json:"version"`
	Type     string       `json:"type"`
	Status   ActiveStatus `json:"status"`
	Progress int          `json:"progress"`
}

func (p *PluginActive) GetKey() string {
	return "PluginActive"
}

func (p *PluginActive) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PluginActive) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *PluginActive) GetField() string {
	return p.Type
}

func init() {
	orm.RegisterModel(&PluginPackage{})
}
