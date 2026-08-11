// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package service
package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
)

const (
	// operationFirstImport 首次导入模式，要求白名单表为空
	operationFirstImport = "firstImport"
	// operationUpdate 覆盖更新模式，事务清表后批量插入
	operationUpdate = "update"
	// maxImportCount 单文件最大导入条数
	maxImportCount = 200000
	// importBatchSize 批量插入分片大小
	importBatchSize = 1000
)

// WhiteListManageService 白名单管理服务
type WhiteListManageService interface {
	ImportIMEIList(reader io.Reader, operation string) (int, error)
	ExportIMEIList() (string, error)
}

type whiteListManageServiceImpl struct {
	whiteListDao whiteListDaoInterface
}

var (
	whiteListManageServiceInstance WhiteListManageService
	whiteListManageServiceOnce     sync.Once
)

// NewWhiteListManageService 获取白名单管理服务单例
func NewWhiteListManageService() WhiteListManageService {
	whiteListManageServiceOnce.Do(func() {
		whiteListManageServiceInstance = &whiteListManageServiceImpl{
			whiteListDao: dao.NewWhiteListDao(),
		}
	})
	return whiteListManageServiceInstance
}

// ImportIMEIList 解析 CSV 强校验后导入白名单，导入全程持写锁
func (s *whiteListManageServiceImpl) ImportIMEIList(reader io.Reader, operation string) (int, error) {
	records, err := parseWhiteListCSV(reader)
	if err != nil {
		return 0, err
	}
	authImportLock.Lock()
	defer authImportLock.Unlock()
	if operation == operationFirstImport {
		count, err := s.whiteListDao.Count()
		if err != nil {
			logger.Errorf("[ImportIMEIList] count white list failed, err: [%v]", err)
			return 0, err
		}
		if count > 0 {
			return 0, errors.New("white list is not empty, please use update")
		}
		if err := s.insertInBatches(records); err != nil {
			return 0, err
		}
		return len(records), nil
	}
	if err := s.whiteListDao.ClearAndInsert(&records); err != nil {
		logger.Errorf("[ImportIMEIList] clear and insert failed, err: [%v]", err)
		return 0, err
	}
	return len(records), nil
}

// ExportIMEIList 导出全量白名单为无 header 的 IMEI/IMSI 两列 CSV 文本
func (s *whiteListManageServiceImpl) ExportIMEIList() (string, error) {
	var list []db.WhiteList
	if err := s.whiteListDao.ListAll(&list); err != nil {
		logger.Errorf("[ExportIMEIList] list white list failed, err: [%v]", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	for _, item := range list {
		if err := writer.Write([]string{item.Imei, item.Imsi}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// parseWhiteListCSV 按无 header 口径解析 CSV，空行自动跳过，逐行强校验 15 位纯数字
func parseWhiteListCSV(reader io.Reader) ([]db.WhiteList, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = 2
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse csv failed: %v", err)
	}
	if len(rows) > maxImportCount {
		return nil, fmt.Errorf("import count %d exceeds limit %d", len(rows), maxImportCount)
	}
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	records := make([]db.WhiteList, 0, len(rows))
	for i, row := range rows {
		if !imeiPattern.MatchString(row[0]) || !imeiPattern.MatchString(row[1]) {
			return nil, fmt.Errorf("invalid imei/imsi format at line %d", i+1)
		}
		records = append(records, db.WhiteList{
			Imei:      row[0],
			Imsi:      row[1],
			CreatedAt: createdAt,
		})
	}
	return records, nil
}

// insertInBatches 按 1000 条/批分片批量插入
func (s *whiteListManageServiceImpl) insertInBatches(records []db.WhiteList) error {
	for start := 0; start < len(records); start += importBatchSize {
		end := start + importBatchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[start:end]
		if err := s.whiteListDao.InsertMulti(&batch); err != nil {
			logger.Errorf("[insertInBatches] insert batch failed, err: [%v]", err)
			return err
		}
	}
	return nil
}
