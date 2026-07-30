// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/models/req"
	"GIDS/models/resp"
	"GIDS/service"
)

// AuthController 终端鉴权与白名单管理入口，仅监听内部服务（127.0.0.1:9090/运维内网段）
type AuthController struct {
	BaseController
	authService   service.AuthService
	manageService service.AuthManageService
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
	c.manageService = service.NewAuthManageService()
}

// AuthIMEI 联合鉴权：格式非法或未命中返回 code=401，均命中返回 code=200
func (c *AuthController) AuthIMEI() {
	request := new(req.AuthIMEIRequest)
	if err := c.RequestBodyUnmarshalTo(request); err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.AuthFailed, Message: "imei or imsi format invalid"})
		return
	}
	pass, formatValid := c.authService.Check(request.IMEI, request.IMSI)
	if !formatValid {
		c.Failed(resp.BaseResponse{Code: retcode.AuthFailed, Message: "imei or imsi format invalid"})
		return
	}
	if !pass {
		c.Failed(resp.BaseResponse{Code: retcode.AuthFailed, Message: "auth rejected"})
		return
	}
	c.OK(resp.BaseResponse{Code: retcode.Success, Message: "success"})
}

// ImportIMEIList 白名单导入：multipart 上传 CSV + form 参数 operation=firstImport/update
func (c *AuthController) ImportIMEIList() {
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		logger.Errorf("[importIMEIList] get form file failed, err: [%v]", err)
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: "form file [file] required"})
		return
	}
	defer file.Close()
	operation := c.Request().FormValue("operation")
	count, errCode, msg := c.manageService.Import(file, operation)
	if errCode != retcode.Success {
		c.Failed(resp.BaseResponse{Code: errCode, Message: msg})
		return
	}
	c.OK(resp.DataResponse{
		BaseResponse: resp.BaseResponse{Code: retcode.Success, Message: "success"},
		Data:         count,
	})
}

// ExportIMEIList 白名单导出：全量白名单生成 CSV 文本返回
func (c *AuthController) ExportIMEIList() {
	csvText, err := c.manageService.Export()
	if err != nil {
		logger.Errorf("[exportIMEIList] export failed, err: [%v]", err)
		c.InternalServiceError()
		return
	}
	c.AddHeader("Content-Type", "text/csv")
	if _, err := c.ResponseWriter().Write([]byte(csvText)); err != nil {
		logger.Errorf("[exportIMEIList] write response failed, err: [%v]", err)
	}
}
