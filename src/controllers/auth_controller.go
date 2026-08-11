// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"net/http"

	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/models/req"
	"GIDS/models/resp"
	"GIDS/service"
)

// maxImportFileSize 导入文件大小硬限 3MB
const maxImportFileSize = 3 * 1024 * 1024

// AuthController 终端鉴权与白名单管理接口，仅注册内网
type AuthController struct {
	BaseController
	authService   service.AuthService
	manageService service.WhiteListManageService
}

func (c *AuthController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/auth/v1/authIMEI":       "POST:AuthIMEI",
			"/auth/v1/importIMEIList": "POST:ImportIMEIList",
			"/auth/v1/exportIMEIList": "GET:ExportIMEIList",
		},
	}
}

func (c *AuthController) Prepare() {
	c.authService = service.NewAuthService()
	c.manageService = service.NewWhiteListManageService()
}

// AuthIMEI 终端联合鉴权，统一 HTTP 200+body code 标识结果
func (c *AuthController) AuthIMEI() {
	request := new(req.AuthIMEIRequest)
	if err := c.RequestBodyUnmarshalTo(request); err != nil {
		c.OK(resp.BaseResponse{Code: retcode.AuthFailed, Message: "format invalid"})
		return
	}
	allowed, formatValid := c.authService.AuthIMEI(request.IMEI, request.IMSI)
	if !formatValid {
		c.OK(resp.BaseResponse{Code: retcode.AuthFailed, Message: "format invalid"})
		return
	}
	if !allowed {
		c.OK(resp.BaseResponse{Code: retcode.AuthFailed, Message: "auth rejected"})
		return
	}
	c.OK(resp.BaseResponse{Code: retcode.Success, Message: "success"})
}

// ImportIMEIList 白名单 CSV 导入
func (c *AuthController) ImportIMEIList() {
	file, header, err := c.GetFile("file")
	if err != nil {
		logger.Errorf("[ImportIMEIList] get file failed, err: [%v]", err)
		c.OK(resp.BaseResponse{Code: retcode.ClientFailed, Message: "invalid parameter"})
		return
	}
	defer file.Close()
	operation := c.GetString("operation")
	if header.Size > maxImportFileSize ||
		(operation != "firstImport" && operation != "update") {
		c.OK(resp.BaseResponse{Code: retcode.ClientFailed, Message: "invalid parameter"})
		return
	}
	count, err := c.manageService.ImportIMEIList(file, operation)
	if err != nil {
		c.OK(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}
	c.OK(resp.DataResponse{
		BaseResponse: resp.BaseResponse{Code: retcode.Success, Message: "success"},
		Data:         count,
	})
}

// ExportIMEIList 白名单 CSV 导出
func (c *AuthController) ExportIMEIList() {
	csvText, err := c.manageService.ExportIMEIList()
	if err != nil {
		c.InternalServiceError()
		return
	}
	c.AddHeader("Content-Type", "text/csv")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	if _, err := c.ResponseWriter().Write([]byte(csvText)); err != nil {
		logger.Errorf("[ExportIMEIList] write response failed, err: [%v]", err)
	}
}
