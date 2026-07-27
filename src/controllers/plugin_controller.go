// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package controllers include web controllers
package controllers

import (
	"errors"
	"runtime"
	"strings"

	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/models/req"
	"GIDS/models/resp"
	"GIDS/service"
)

type PluginController struct {
	BaseController
	pluginService  service.PluginService
	browserService service.BrowserService
}

func (c *PluginController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/plugin/v1/upload":  "POST:UploadPluginPackage",
			"/plugin/v1/delete":  "POST:DeletePluginPackage",
			"/plugin/v1/getAll":  "POST:GetPluginPackages",
			"/plugin/v1/load":    "POST:LoadPlugin",
			"/plugin/v1/current": "POST:GetCurrentPlugins",
		},
	}
}

func (c *PluginController) Prepare() {
	c.pluginService = service.NewPluginService()
	c.browserService = service.NewBrowserService()
}

func getStackTrace() string {
	buf := make([]byte, 1024*10) // 10KB的缓冲区，足够存储一般堆栈信息
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func (c *PluginController) UploadPluginPackage() {
	// 使用defer和recover捕获异常
	defer func() {
		if err := recover(); err != nil {
			// 捕获到panic，进行处理
			logger.Errorf("UploadPluginPackage error: %v", err)
			stackTrace := getStackTrace()
			// 格式化输出，每行前缀添加"  "
			formattedTrace := strings.ReplaceAll(stackTrace, "\n", "\n  ")
			logger.Errorf("  %s\n", formattedTrace)
		}
	}()
	request, err := c.parseUploadPluginPackageParam()
	if err != nil {
		c.Failed(resp.BaseResponse{
			Code:    retcode.ClientFailed,
			Message: "param parse failed",
		})
		return
	}

	err = c.pluginService.UploadPluginPackage(&request)
	if err != nil {
		logger.Errorf("upload plugin package failed! err is %v", err)
		c.InternalServiceError()
		return
	}
	c.OK(nil)
}

func (c *PluginController) parseUploadPluginPackageParam() (req.UploadPluginPackageReq, error) {
	filename := c.Request().FormValue("filename")
	file, header, err := c.Request().FormFile("file")
	logger.Infof("ncx: filename is: %s", filename)
	if err != nil {
		logger.Errorf("form file error: %v", err)
		return req.UploadPluginPackageReq{}, err
	}
	if header == nil {
		logger.Errorf("form file error: header is nil")
		return req.UploadPluginPackageReq{}, errors.New("file header is nil")
	}

	var request = req.UploadPluginPackageReq{
		Filename: filename,
		File:     file,
		Size:     header.Size,
	}
	if filename == "" {
		request.Filename = header.Filename
	}
	err = request.Validate()
	if err != nil {
		logger.Errorf("form file error: %v", err)
		return req.UploadPluginPackageReq{}, err
	}
	return request, err
}

func (c *PluginController) DeletePluginPackage() {
	var request req.PluginPackageReq
	err := c.RequestBodyUnmarshalTo(&request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}
	err = c.pluginService.DeletePluginPackage(&request)
	if err != nil {
		c.InternalServiceError()
		return
	}
	c.OK(nil)
}

func (c *PluginController) GetPluginPackages() {
	pluginInfos, err := c.pluginService.GetPluginPackages()
	if err != nil {
		logger.Errorf("get plugin info failed! err is %v", err)
		c.InternalServiceError()
		return
	}
	c.OK(pluginInfos)
}

func (c *PluginController) LoadPlugin() {
	var request req.PluginPackageReq
	err := c.RequestBodyUnmarshalTo(&request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}
	instances := c.browserService.GetAllServiceInstances()

	err = c.pluginService.LoadPlugin(&request, instances)
	if err != nil {
		c.InternalServiceError()
		return
	}
	c.OK(nil)
}

func (c *PluginController) GetCurrentPlugins() {
	pluginInfos, err := c.pluginService.GetCurrentPlugins()
	if err != nil {
		c.InternalServiceError()
		return
	}
	c.OK(pluginInfos)
}
