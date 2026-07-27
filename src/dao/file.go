// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package dao

import (
	"GIDS/models/db"
)

type FileDao struct {
	BaseInterface
}

func (f *FileDao) Exist(bucket, name string) (bool, error) {
	var cnt int
	err := f.QueryOne(&cnt, "select count(*) from t_file where bucket = ? and name = ?", bucket, name)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func NewFileDao() *FileDao {
	dao := &FileDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.File{},
	}
	return dao
}
