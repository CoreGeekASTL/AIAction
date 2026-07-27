/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package dao
package dao

import "GIDS/models/db"

// ConfigCenterDao config center dao
type ConfigCenterDao struct {
	BaseInterface
}

// NewConfigCenterDao create new dao
func NewConfigCenterDao() *ConfigCenterDao {
	dao := &ConfigCenterDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.ConfigCenter{},
	}
	return dao
}
