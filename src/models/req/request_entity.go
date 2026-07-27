// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package req
package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type ITable interface {
	IRequest
	orm.TableNameI
}

type IRequest interface {
	Validate() error
}

type UserIdentity struct {
	IMSI string `json:"imsi"`
	IMEI string `json:"imei"`
}

type LoginAuthRequest struct {
	UserIdentity
	Manufacturer   string `json:"manufacturer"`
	Model          string `json:"model"`
	AppType        string `json:"appType"`
	ExtendModel    string `json:"extendModel"`
	Country        string `json:"country"`
	Platform       string `json:"platform"`
	Width          string `json:"width"`
	Height         string `json:"height"`
	MCC            string `json:"mcc"`
	MNC            string `json:"mnc"`
	Lac            string `json:"lac"`
	CI             string `json:"ci"`
	Rxlev          string `json:"rxlev"`
	TotalKb        string `json:"totalKb"`
	FreeKb         string `json:"freeKb"`
	ClientLanguage string `json:"clientLanguage"`
	DeviceType     string `json:"deviceType"`
}

func (g LoginAuthRequest) Validate() error {
	return nil
}

type SyncBrowserConfigRequest struct {
	RouteAPPConfigList []RouteAppConfig `json:"routeAppConfigList,omitempty"`
	ChromeConfigList   []ChromeConfig   `json:"chromeConfigList,omitempty"`
	URLConfigs         []URLConfigs     `json:"urlConfigList,omitempty"`
}

type RouteAppConfig struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	AppType      int    `json:"type,omitempty"`
	Mode         int    `json:"mode,omitempty"`
	ExtendModel  string `json:"extendModel,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
}

type ChromeConfig struct {
	Manufacturer   string `json:"manufacturer,omitempty"`
	Model          string `json:"model,omitempty"`
	Country        string `json:"country,omitempty"`
	AppFrameRate   int    `json:"appFrameRate,omitempty"`
	VideoFrameRate int    `json:"videoFrameRate,omitempty"`
	AppBitRate     int    `json:"appBitRate,omitempty"`
	VideoBitRate   int    `json:"videoBitRate,omitempty"`
	SampleRate     int    `json:"sampleRate,omitempty"`
	Channels       int    `json:"channels,omitempty"`
	MachineType    int    `json:"machineType,omitempty"`
	FFCode         string `json:"ffCode,omitempty"`
	Resolution     string `json:"resolution,omitempty"`
	RecordMode     int    `json:"recordMode,omitempty"`
}

type URLConfigs struct {
	NodeIdent   string `json:"nodeIdent,omitempty"`
	APPType     int    `json:"appType,omitempty"`
	URL         string `json:"url,omitempty"`
	AppID       string `json:"appID,omitempty"`
	Name        string `json:"name,omitempty"`
	IsVideoType bool   `json:"isVideoType,omitempty"`
	IsWebType   bool   `json:"isWebType,omitempty"`
	IsShortType bool   `json:"isShortType,omitempty"`
	UserAgent   string `json:"userAgent,omitempty"`
}

func (g SyncBrowserConfigRequest) Validate() error {
	return nil
}

type DeleteCacheRequest struct {
	IMEI string `form:"imei" json:"imei" binding:"required"`
	IMSI string `form:"imsi" json:"imsi" binding:"required"`
}

func (r DeleteCacheRequest) Validate() error {
	var errors []string

	if r.IMEI == "" {
		errors = append(errors, "IMEI cannot be empty")
	}

	if r.IMSI == "" {
		errors = append(errors, "IMSI cannot be empty")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}
	return nil
}

type FileUploadRequest struct {
	FileName string `form:"fileName" binding:"required"`
}

func (r FileUploadRequest) Validate() error {
	if r.FileName == "" {
		return errors.New("fileName cannot be empty")
	}
	return nil
}

type UpdateUserBindRequest struct {
	SessionID            string `json:"sessionID"`
	BrowserInstance      string `json:"browserInstance"`
	MediaEndpoint        string `json:"mediaEndpoint"`
	ControlEndpoint      string `json:"controlEndpoint"`
	MediaTlsEndpoint     string `json:"mediaTlsEndpoint"`
	ControlTlsEndpoint   string `json:"controlTlsEndpoint"`
	InnerMediaEndpoint   string `json:"innerMediaEndpoint"`
	InnerBrowserEndpoint string `json:"innerBrowserEndpoint"`
}

func (r UpdateUserBindRequest) Validate() error {
	if r.SessionID == "" {
		return errors.New("sessionID cannot be empty")
	}
	return nil
}

// 上报多条数据请求
type MultiTableRequest struct {
	Items []json.RawMessage `json:"items"`
}

func (m *MultiTableRequest) Validate() error {
	if m == nil {
		return errors.New("multiTableRequest is nil")
	}
	if len(m.Items) == 0 {
		return errors.New("items cannot be empty")
	}
	return nil
}
