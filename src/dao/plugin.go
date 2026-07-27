// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package dao

import "GIDS/models/db"

type PluginPackageDao struct {
	BaseInterface
}

func NewPluginPackageDao() *PluginPackageDao {
	dao := &PluginPackageDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.PluginPackage{},
	}
	return dao
}
