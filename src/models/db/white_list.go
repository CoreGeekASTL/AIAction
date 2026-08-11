// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package db
package db

import (
	"github.com/beego/beego/v2/client/orm"
)

// WhiteList 白名单实体，对应 t_white_list 表，IMEI+IMSI 精确匹配组合
type WhiteList struct {
	Imei      string `orm:"pk;column(imei);size(15)"`
	Imsi      string `orm:"column(imsi);size(15)"`
	CreatedAt string `orm:"column(created_at)"`
}

func (w *WhiteList) TableName() string {
	return "t_white_list"
}

func init() {
	orm.RegisterModel(&WhiteList{})
}
