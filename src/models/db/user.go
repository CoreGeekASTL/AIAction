// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package db
package db

import (
	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Key          string `orm:"pk;column(key)"`
	Manufacturer string `json:"manufacturer" orm:"column(manufacturer)"`
	Model        string `json:"model" orm:"column(model)"`
	ExtendModel  string `json:"extendModel" orm:"column(extend_model)"`
	Country      string `json:"country" orm:"column(country)"`
	Platform     string `json:"platform" orm:"column(platform)"`
	Width        string `json:"width" orm:"column(width)"`
	Height       string `json:"height" orm:"column(height)"`
	MCC          string `json:"mcc" orm:"column(mcc)"`
	MNC          string `json:"mnc" orm:"column(mnc)"`
	DeviceType   string `json:"deviceType" orm:"column(device_type)"`
	CreatedAt    string `json:"-" orm:"column(created_at)"`
	UpdatedAt    string `json:"-" orm:"column(updated_at)"`
}

func (u *User) TableName() string {
	return "t_user"
}

func (u *User) GetKey() string {
	return u.Key
}

type UserBind struct {
	Key                  string `json:"-" orm:"pk;column(key)"`
	BrowserInstance      string `json:"browserInstance" orm:"column(browser_instance)"`
	BrowserCap           int    `json:"-" orm:"-"`
	MediaEndpoint        string `json:"mediaEndpoint" orm:"column(media_endpoint)"`
	ControlEndpoint      string `json:"controlEndpoint" orm:"column(control_endpoint)"`
	MediaTlsEndpoint     string `json:"mediaTlsEndpoint" orm:"column(media_tls_endpoint)"`
	ControlTlsEndpoint   string `json:"controlTlsEndpoint" orm:"column(control_tls_endpoint)"`
	InnerMediaEndpoint   string `json:"innerMediaEndpoint" orm:"column(inner_media_endpoint)"`
	InnerBrowserEndpoint string `json:"innerBrowserEndpoint" orm:"column(inner_browser_endpoint)"`
	Token                string `json:"token" orm:"column(token)"`
	Heartbeats           string `json:"heartbeats" orm:"column(updated_at)"`
}

func (u *UserBind) GetField() string {
	return u.Key
}

func (u *UserBind) TableName() string {
	return "t_user_bind"
}

func init() {
	orm.RegisterModel(&User{})
	orm.RegisterModel(&UserBind{})
}
