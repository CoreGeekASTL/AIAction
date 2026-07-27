// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package resp
package resp

import "GIDS/models/db"

type PluginPackageResponse struct {
	BaseResponse   `json:"baseResponse"`
	PluginPackages []db.PluginPackage `json:"data"`
}

type PluginPackage struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Type        string `json:"type"`
	PackageName string `json:"packageName"`
}

type PluginInfoResponse struct {
	BaseResponse `json:"baseResponse"`
	PluginInfos  []PluginInfo `json:"data"`
}

type PluginInfo struct {
	Name     string          `json:"name"`
	Version  string          `json:"version"`
	Type     string          `json:"type"`
	Status   db.ActiveStatus `json:"status"`
	Progress int             `json:"progress"`
}
