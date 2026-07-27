// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package dao

import "GIDS/models/db"

type UserDao struct {
	BaseInterface
}

func NewUserDaoDao() *UserDao {
	dao := &UserDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.User{},
	}
	return dao
}

type UserBindDao struct {
	BaseInterface
}

func NewUserBindDao() *UserBindDao {
	dao := &UserBindDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.UserBind{},
	}
	return dao
}
