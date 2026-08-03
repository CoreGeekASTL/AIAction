// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
)

const (
	// importMaxFileSize 导入文件大小上限 3MB
	importMaxFileSize = 3 * 1024 * 1024
	// importMaxRecords 单文件记录数上限 20W
	importMaxRecords = 200000
	// csvFieldCount CSV 每行字段数：IMEI,IMSI 两列
	csvFieldCount = 2
	// timeLayoutSecond 秒级时间格式，CreatedAt 落库用
	timeLayoutSecond = "2006-01-02 15:04:05"
	// operationFirstImport 首次导入，要求白名单表为空
	operationFirstImport = "firstImport"
	// operationUpdate 覆盖更新，事务清表 + 批量插入
	operationUpdate = "update"
)

// AuthManageService 白名单管理服务：CSV 解析校验 + 批量入库 + 导出生成，不做鉴权
type AuthManageService interface {
	// Import 导入 CSV，返回导入条数、错误码（200/-1/-2）与消息
	Import(reader io.Reader, operation string) (count int, errCode int, msg string)
	// Export 导出全量白名单为 CSV 文本（带 IMEI,IMSI 表头）
	Export() (string, error)
}

type authManageServiceImpl struct {
	store whiteListDao
	auth  AuthService
}

var _ AuthManageService = &authManageServiceImpl{}

var (
	authManageOnce     sync.Once
	authManageInstance *authManageServiceImpl
)

func NewAuthManageService() AuthManageService {
	authManageOnce.Do(func() {
		authManageInstance = &authManageServiceImpl{
			store: dao.NewWhiteListDao(),
			auth:  NewAuthService(),
		}
	})
	return authManageInstance
}

// Import 解析校验 CSV 并批量入库；任一校验失败整批不加载；事务提交成功后清空鉴权缓存立即生效
func (s *authManageServiceImpl) Import(reader io.Reader, operation string) (int, int, string) {
	if operation != operationFirstImport && operation != operationUpdate {
		return 0, retcode.ClientFailed, fmt.Sprintf("operation [%s] invalid, expect firstImport or update", operation)
	}
	content, errCode, errMsg := readWithSizeLimit(reader)
	if errCode != retcode.Success {
		return 0, errCode, errMsg
	}
	records, parseMsg := parseWhiteListCSV(content)
	if parseMsg != "" {
		return 0, retcode.InternalFailed, parseMsg
	}
	if operation == operationFirstImport {
		count, err := s.store.Count()
		if err != nil {
			logger.Errorf("[importIMEIList] count white list failed, err: [%v]", err)
			return 0, retcode.InternalFailed, "count white list failed"
		}
		if count != 0 {
			return 0, retcode.InternalFailed, "white list table not empty, please use update operation"
		}
		if err := s.store.InsertMulti(records); err != nil {
			logger.Errorf("[importIMEIList] insert multi failed, err: [%v]", err)
			return 0, retcode.InternalFailed, "insert white list failed"
		}
	} else {
		if err := s.store.ClearAndInsert(records); err != nil {
			logger.Errorf("[importIMEIList] clear and insert failed, err: [%v]", err)
			return 0, retcode.InternalFailed, "update white list failed"
		}
	}
	s.auth.ClearCache()
	return len(records), retcode.Success, "success"
}

// readWithSizeLimit 读取全部内容并校验大小上限，返回错误码与消息（Success 表示通过）
func readWithSizeLimit(reader io.Reader) ([]byte, int, string) {
	content, err := io.ReadAll(reader)
	if err != nil {
		logger.Errorf("[importIMEIList] read csv failed, err: [%v]", err)
		return nil, retcode.InternalFailed, "read csv file failed"
	}
	if len(content) > importMaxFileSize {
		return nil, retcode.InternalFailed, "csv file size exceeds 3MB"
	}
	return content, retcode.Success, ""
}

// Export 全量查询白名单生成 CSV 文本，首行为 IMEI,IMSI 表头
func (s *authManageServiceImpl) Export() (string, error) {
	list, err := s.store.ListAll()
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	builder.WriteString("IMEI,IMSI\n")
	for _, item := range list {
		writeCSVRow(&builder, item)
	}
	return builder.String(), nil
}

// writeCSVRow 追加一行 IMEI,IMSI 数据行
func writeCSVRow(builder *strings.Builder, item db.AuthWhitelist) {
	builder.WriteString(item.IMEI)
	builder.WriteString(",")
	builder.WriteString(item.IMSI)
	builder.WriteString("\n")
}

// parseWhiteListCSV 纯数据 CSV（无 header），所有行当数据逐行强校验；返回空错误串表示成功
func parseWhiteListCSV(content []byte) ([]db.AuthWhitelist, string) {
	csvReader := csv.NewReader(strings.NewReader(string(content)))
	csvReader.FieldsPerRecord = csvFieldCount
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Sprintf("parse csv failed: %v", err)
	}
	if len(rows) == 0 {
		return nil, "csv file is empty"
	}
	if len(rows) > importMaxRecords {
		return nil, fmt.Sprintf("csv records %d exceed max %d", len(rows), importMaxRecords)
	}
	records := make([]db.AuthWhitelist, 0, len(rows))
	seen := make(map[string]struct{}, len(rows))
	now := time.Now().Format(timeLayoutSecond)
	for i, row := range rows {
		imei := strings.TrimSpace(row[0])
		imsi := strings.TrimSpace(row[1])
		if !idPattern.MatchString(imei) || !idPattern.MatchString(imsi) {
			return nil, fmt.Sprintf("line %d format invalid: imei or imsi is not 15 digits", i+1)
		}
		key := imei + "_" + imsi
		if _, ok := seen[key]; ok {
			return nil, fmt.Sprintf("line %d duplicate imei+imsi [%s]", i+1, imei)
		}
		seen[key] = struct{}{}
		records = append(records, db.AuthWhitelist{IMEI: imei, IMSI: imsi, CreatedAt: now})
	}
	return records, ""
}
