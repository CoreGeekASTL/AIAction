/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package controllers

import (
	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/models/resp"
)

type TestController struct {
	BaseController
}

func (c *TestController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/test/v1/get": "GET:GetData",
		},
	}
}

func (c *TestController) GetData() {
	logger.Infof("test controller get request")
	response := resp.DataResponse{
		BaseResponse: resp.BaseResponse{
			Code:    retcode.Success,
			Message: "test success",
		},
		Data: true,
	}
	c.OK(response)
}

func (c *TestController) Prepare() {
	logger.Infof("prepare successful")
}
