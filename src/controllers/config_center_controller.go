/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package controllers
package controllers

import (
	"GIDS/common/constants/retcode"
	"GIDS/models/db"
	"GIDS/models/resp"
	"GIDS/service"
)

// ConfigCenterController config center service controller
type ConfigCenterController struct {
	BaseController
}

// RouteInfo route
func (c *ConfigCenterController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/configCenter/v1/":    "POST:InsertOrUpdate",
			"/configCenter/v1/get": "POST:GetFromDB",
		},
	}
}

// GetFromDB get config from db
func (c *ConfigCenterController) GetFromDB() {
	var request db.ConfigCenter
	err := c.RequestBodyUnmarshalTo(&request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}
	if request.Key == "" {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: "key cannot be empty"})
		return
	}
	fromDB, _ := service.NewConfigCenterService().GetFromDB(request.Key)
	c.OK(fromDB)
}

// InsertOrUpdate insert or update config
func (c *ConfigCenterController) InsertOrUpdate() {
	var request db.ConfigCenter
	err := c.RequestBodyUnmarshalTo(&request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}
	if request.Key == "" {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: "key cannot be empty"})
		return
	}
	err = service.NewConfigCenterService().InsertOrUpdateConfig(request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}
	c.OK(nil)
}
