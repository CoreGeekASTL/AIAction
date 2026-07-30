// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package db
package db

import (
	"github.com/beego/beego/v2/client/orm"
)

// AuthWhitelist 终端鉴权白名单，IMEI+IMSI 联合鉴权，IMEI 为主键，IMSI 唯一性由 DDL 联合索引兜底
type AuthWhitelist struct {
	IMEI      string `json:"imei" orm:"pk;column(imei)"`
	IMSI      string `json:"imsi" orm:"column(imsi)"`
	CreatedAt string `json:"-" orm:"column(created_at)"`
}

func (w *AuthWhitelist) TableName() string {
	return "t_white_list"
}

func init() {
	orm.RegisterModel(&AuthWhitelist{})
}
