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

type FileController struct {
	BaseController
	fileService service.FileService
}

func (c *FileController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/app-api/control/file/upload":             "POST:HandleUpload",
			"/app-api/control/file/download/:fileName": "GET:HandleDownload",

			"/file/v1/:bucket/:name":         "GET:Download",
			"/file/v1/:bucketName/:name":     "POST:Upload",
			"/file/v1/:bucket/:name/exist":   "GET:Exist",
			"/file/v1/:bucketName/:fileName": "DELETE:HandleDelete",
		},
	}
}

func (c *FileController) Exist() {
	bucket := c.PathParameter(":bucket")
	name := c.PathParameter(":name")
	ok, err := c.fileService.Exist(bucket, name)
	if err != nil {
		logger.Errorf("file exist failed: %v", err)
		c.InternalServiceError()
		return
	}
	if !ok {
		c.NotFound()
		return
	}
	c.OK(nil)
}

func (c *FileController) Upload() {
	bucket := c.PathParameter(":bucketName")
	name := c.PathParameter(":name")
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		logger.Errorf("form file error: %v", err)
		c.InternalServiceError()
		return
	}
	if header == nil {
		logger.Errorf("form file error: header is nil")
		c.InternalServiceError()
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Errorf("read file error: %v", err)
		c.InternalServiceError()
		return
	}
	filePath, err := c.fileService.UploadFile(bucket, name, content)
	if err != nil {
		logger.Errorf("upload file error: %v", err)
		c.InternalServiceError()
		return
	}
	c.OK(filePath)
}

func (c *FileController) Download() {
	bucket := c.PathParameter(":bucket")
	name := c.PathParameter(":name")
	content, err := c.fileService.DownloadFile(bucket, name)
	if err != nil {
		logger.Errorf("download file failed: %v", err)
		c.InternalServiceError()
		return
	}
	c.AddHeader("Content-Disposition", "attachment; filename="+name)
	c.AddHeader("Content-Type", "application/octet-stream")
	c.AddHeader("Content-Length", string(rune(len(content))))

	if _, err := c.ResponseWriter().Write(content); err != nil {
		c.InternalServiceError()
		return
	}
}

func (c *FileController) Prepare() {
	c.fileService = service.NewFileService()
}

func (c *FileController) HandleUpload() {
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

func (c *FileController) HandleDownload() {
	fileName := c.PathParameter(":fileName")
	if fileName == "" {
		c.Failed(resp.BaseResponse{
			Code:    retcode.InternalFailed,
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

func (c *FileController) HandleDelete() {
	bucket := c.PathParameter(":bucketName")
	name := c.PathParameter(":fileName")
	err := c.fileService.DeleteFile(bucket, name)
	if err != nil {
		logger.Errorf("delete file failed: %v", err)
		c.InternalServiceError()
		return
	}
	c.AddHeader("Content-Disposition", "attachment; filename="+name)
	c.AddHeader("Content-Type", "application/octet-stream")
	c.OK(nil)
}
