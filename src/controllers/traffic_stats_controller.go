// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
	"GIDS/models/req"
	"GIDS/models/resp"
	"GIDS/service"
	util "GIDS/utils/fileutil"
)

const (
	SessionStatsCsv = "session_stats.csv"
	MediaStatsCsv   = "media_stats.csv"
	ControlStatsCsv = "control_stats.csv"
	MonthFormat     = "2006-01"
)

type TrafficStatsController struct {
	BaseController
	mts                 *dao.MediaTrafficStatsDao
	cts                 *dao.ControlTrafficStatsDao
	sl                  *dao.SessionStatsDao
	trafficStatsService service.TrafficStatsService
}

func (c *TrafficStatsController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/stats/v1/session":                 "POST:SessionStats",
			"/stats/v1/traffic/media":           "POST:MediaTrafficStats",
			"/stats/v1/traffic/control":         "POST:ControlTrafficStats",
			"/stats/v1/exportStaticData/:month": "GET:ExportStaticData",
		},
	}
}

func (c *TrafficStatsController) Prepare() {
	c.mts = dao.NewMediaTrafficStatsDao()
	c.cts = dao.NewControlTrafficStatsDao()
	c.sl = dao.NewSessionLogDao()
	c.trafficStatsService = service.NewTrafficStatsService("")
}

func (c *TrafficStatsController) SessionStats() {
	c.insertOrUpdateSessionStats()
}

func (c *TrafficStatsController) MediaTrafficStats() {
	c.insertMultiData("media traffic stats")
}

func (c *TrafficStatsController) ControlTrafficStats() {
	c.insertMultiData("control traffic stats")
}

func (c *TrafficStatsController) insertOrUpdateSessionStats() {
	// 创建参数实例
	session := new(db.SessionStats)

	// 解析请求体
	err := c.RequestBodyUnmarshalTo(session)
	if err != nil {
		logger.Errorf("[session stats] unmarshal failed, err: [%v]", err)
		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}

	// 调用Service层处理
	logger.Infof("[session stats] insert Or Update SessionStats")
	if err := c.trafficStatsService.HandleSessionStats(session); err != nil {
		logger.Errorf("[session stats] handle session failed, err: [%v]", err)
		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}

	c.OK(nil)
}

func (c *TrafficStatsController) insertMultiData(tag string) {
	// 创建参数实例
	var multiTableReq req.MultiTableRequest

	// 解析请求体
	err := c.RequestBodyUnmarshalTo(&multiTableReq)
	if err != nil {
		logger.Errorf("[%v] unmarshal failed, err: [%v]", tag, err)
		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}

	// 调用Service层处理
	logger.Infof("[stats] Batch Insert Stats")
	err = c.trafficStatsService.BatchInsertStats(tag, multiTableReq.Items)
	if err != nil {
		logger.Errorf("[%v] batch insert failed, err is [%v]", tag, err)
		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}
	logger.Infof("batch insert success")

	c.OK(nil)
}

func (c *TrafficStatsController) ExportStaticData() {
	logger.Infof("begin to export static data")

	if err := c.validateExportStaticDataParam(); err != nil {
		logger.Errorf("validate failed : %v", err)
		c.Failed(resp.BaseResponse{
			Code:    retcode.ClientFailed,
			Message: err.Error(),
		})
		return
	}
	// 创建临时目录
	month := c.PathParameter(":month")
	tmpDir, err := os.MkdirTemp("", "export_*")
	if err != nil {
		logger.Errorf("create temp dir failed, err: %v", err)
		c.InternalServiceError()
		return
	}
	defer os.RemoveAll(tmpDir)
	// 查询数据并写入 csv文件
	if err = c.exportCSVFiles(tmpDir, month); err != nil {
		logger.Errorf("failed to export CSV files, err: %v", err)
		c.InternalServiceError()
		return
	}
	// 压缩文件
	zipPath := filepath.Join(tmpDir, month+".zip")
	if err = util.CreateZipFile(zipPath, tmpDir); err != nil {
		logger.Errorf("failed to create zip file, err: %v", err)
		c.InternalServiceError()
		return
	}
	// 响应返回zip包
	if err = sendZipResponse(c, month, zipPath); err != nil {
		logger.Errorf("failed to write response, err: %v", err)
		c.InternalServiceError()
		return
	}

	logger.Infof("Success to export static data")
}

func (c *TrafficStatsController) exportCSVFiles(tmpDir, month string) error {
	csvFiles := []struct {
		name       string
		exportFunc func(string, string) error
	}{
		{SessionStatsCsv, c.trafficStatsService.ExportSessionList},
		{MediaStatsCsv, c.trafficStatsService.ExportMediaList},
		{ControlStatsCsv, c.trafficStatsService.ExportControlList},
	}

	for _, file := range csvFiles {
		filePath := filepath.Join(tmpDir, file.name)
		if err := file.exportFunc(month, filePath); err != nil {
			return fmt.Errorf("failed to export %s data, err: %v", file.name, err)
		}
	}
	return nil
}

func (c *TrafficStatsController) validateExportStaticDataParam() error {
	month := c.PathParameter(":month")
	if month == "" {
		return fmt.Errorf("param month is required")
	}

	_, err := time.Parse(MonthFormat, month)
	if err != nil {
		return fmt.Errorf("invalid month value : %v", err)
	}
	return nil
}

func sendZipResponse(c *TrafficStatsController, month, zipPath string) error {
	c.AddHeader("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", month))
	c.AddHeader("Content-Type", "application/zip")

	zipFile, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	_, err = io.Copy(c.ResponseWriter(), zipFile)
	return err
}
