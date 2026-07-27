// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package browsergateway
package browsergateway

type ExtensionLoadRequest struct {
	BucketName        string `json:"bucket_name"`
	ExtensionFilePath string `json:"extension_file_path"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Type              string `json:"type"`
}

type InitBrowserRequest struct {
	Factory        string `json:"factory"`
	DevType        string `json:"dev_type"`
	ExType         string `json:"ext_type"`
	PlatType       string `json:"plat_type"`
	LcdWidth       string `json:"lcd_width"`
	LcdHeight      string `json:"lcd_height"`
	AppType        string `json:"app_type"`
	AppID          string `json:"appid"`
	IMSI           string `json:"imsi"`
	IMEI           string `json:"imei"`
	DeviceType     string `json:"device_type"`
	ClientLanguage string `json:"client_language"`
	PlayMode       string `json:"play_mode"`
}
