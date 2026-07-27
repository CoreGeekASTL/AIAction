// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/models/req"
	"GIDS/models/resp"
	"GIDS/service"
)

type CacheController struct {
	BaseController
}

func (c *CacheController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/app-api/devicetcp/cache/deleteCache": "POST:DeleteCache",
		},
	}
}

func (c *CacheController) Prepare() {

}

func (c *CacheController) DeleteCache() {
	request := new(req.DeleteCacheRequest)
	err := c.RequestBodyUnmarshalTo(request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}

	err = service.DeleteCache(request.IMEI, request.IMSI)
	if err != nil {
		logger.Errorf("delete cache failed, err is %v", err)
		para := &logger.AuditsPara{
			OperationZH: "删除用户数据",
			OperationEN: "Delete User Data",
			OperateType: logger.DELETE,
			Level:       logger.MinorLevel,
			Username:    "",
			Terminal:    "",
			Result:      1,
			Detail:      "delete failed",
			DetailZH:    "删除失败",
		}
		logger.AuditsLog(para, logger.OpsLog)

		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}

	response := resp.DataResponse{
		BaseResponse: resp.BaseResponse{
			Code:    retcode.Success,
			Message: "Delete success",
		},
		Data: true,
	}
	c.OK(response)
}
