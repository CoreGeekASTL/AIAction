// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package controllers include web controllers
package controllers

import (
	"io"

	"GIDS/common/constants"
	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/models/resp"
	"GIDS/service"
)

// ExFileController for external communication, use https
type ExFileController struct {
	BaseController
	fileService service.FileService
}

func (c *ExFileController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/app-api/control/file/upload":             "POST:HandleUpload",
			"/app-api/control/file/download/:fileName": "GET:HandleDownload",
		},
	}
}

func (c *ExFileController) Prepare() {
	c.fileService = service.NewFileService()
}

func (c *ExFileController) HandleUpload() {
	// 从查询参数获取文件名
	fileName := c.QueryParameter("fileName")
	if fileName == "" {
		c.Failed(resp.BaseResponse{
			Code:    retcode.ClientFailed,
			Message: "Missing fileName parameter",
		})
		return
	}

	// 读取请求体中的文件内容
	fileContent, err := io.ReadAll(c.Body())
	if err != nil {
		logger.Errorf("read file content error: %v", err)
		c.Failed(resp.BaseResponse{
			Code:    retcode.InternalFailed,
			Message: "Failed to read file content",
		})
		return
	}

	// 调用Service层处理文件上传
	path, err := c.fileService.UploadFile(constants.UploadBucket, fileName, fileContent)
	if err != nil {
		logger.Errorf("upload file failed: %v", err)
		c.Failed(resp.BaseResponse{
			Code:    retcode.InternalFailed,
			Message: err.Error(),
		})
		return
	}

	c.OK(resp.DataResponse{
		BaseResponse: resp.BaseResponse{
			Code:    retcode.Success,
			Message: "Files upload success",
		},
		Data: path,
	})
}

func (c *ExFileController) HandleDownload() {
	fileName := c.PathParameter(":fileName")
	if fileName == "" {
		c.Failed(resp.BaseResponse{
			Code:    retcode.ClientFailed,
			Message: "Missing fileName parameter",
		})
		return
	}

	// 调用Service层下载文件
	file, err := c.fileService.DownloadFile(constants.UploadBucket, fileName)
	if err != nil {
		logger.Errorf("download file failed: %v", err)
		c.Failed(resp.BaseResponse{
			Code:    retcode.InternalFailed,
			Message: err.Error(),
		})
		return
	}

	// 设置响应头
	c.AddHeader("Content-Disposition", "attachment; filename="+fileName)
	c.AddHeader("Content-Type", "application/octet-stream")

	// 将文件内容写入响应体
	if _, err := c.ResponseWriter().Write(file); err != nil {
		logger.Errorf("write file to response failed: %v", err)
		c.Failed(resp.BaseResponse{
			Code:    retcode.InternalFailed,
			Message: "Failed to write file to response",
		})
		return
	}
}
