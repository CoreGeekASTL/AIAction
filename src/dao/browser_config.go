// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package dao

import "GIDS/models/db"

type ConfigDao struct {
	BaseInterface
}

func NewConfigDao() *ConfigDao {
	dao := &ConfigDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.Config{},
	}
	return dao
}
