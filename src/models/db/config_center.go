/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package db
package db

import "github.com/beego/beego/v2/client/orm"

// ConfigCenter config enter entity
type ConfigCenter struct {
	ID        int    `json:"-" orm:"auto;pk;column(id)"`
	Key       string `json:"key" orm:"column(config_key)"`
	Value     string `json:"value" orm:"column(config_value)"`
	Describe  string `json:"describe" orm:"column(config_describe)"`
	Enable    bool   `json:"enable" orm:"column(enable)"`
	UpdatedAt string `orm:"column(updated_at)"`
}

// TableName get table name
func (c *ConfigCenter) TableName() string {
	return "t_config_center"
}

func init() {
	orm.RegisterModel(&ConfigCenter{})
}

// Validate IRequest implement
func (c *ConfigCenter) Validate() error {
	return nil
}
