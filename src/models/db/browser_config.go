// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

package db

import (
	"github.com/beego/beego/v2/client/orm"
)

type Config struct {
	ID        int    `orm:"auto;pk;column(id)"`
	Type      string `orm:"column(type)"`
	Content   string `orm:"column(content);type(text)"`
	CreatedAt string `orm:"column(created_at)"`
	UpdatedAt string `orm:"column(updated_at)"`
}

func (c *Config) TableName() string {
	return "t_config"
}

type RouterAPPConfig struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Type         int    `json:"type"`
	Mode         int    `json:"mode"`
	ExtendModel  string `json:"extendModel"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type ChromeConfig struct {
	Manufacturer   string `json:"manufacturer"`
	Model          string `json:"model"`
	Country        string `json:"country"`
	AppFrameRate   int    `json:"appFrameRate"`
	VideoFrameRate int    `json:"videoFrameRate"`
	AppBitRate     int    `json:"appBitRate"`
	VideoBitRate   int    `json:"videoBitRate"`
	SampleRate     int    `json:"sampleRate"`
	Channels       int    `json:"channels"`
	MachineType    int    `json:"machineType"`
	FFCode         string `json:"ffCode"`
	Resolution     string `json:"resolution"`
	RecordMode     int    `json:"recordMode"`
}

type URLConfig struct {
	NodeIdent   string `json:"nodeIdent"`
	APPType     int    `json:"appType"`
	URL         string `json:"url"`
	AppID       string `json:"appID"`
	Name        string `json:"name"`
	IsVideoType bool   `json:"isVideoType"`
	IsWebType   bool   `json:"isWebType"`
	IsShortType bool   `json:"isShortType"`
	UserAgent   string `json:"userAgent"`
}

func init() {
	orm.RegisterModel(&Config{})
}
