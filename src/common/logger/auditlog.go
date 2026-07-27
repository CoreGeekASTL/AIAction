/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package logger

import (
	"GIDS/common/constants"
	gsfapi "Go-chassis-extend/api/GSF/api"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	OpsLog = "cse://AuditLog/plat/audit/v1/logs"
	SecLog = "cse://AuditLog/plat/audit/v1/seculogs"
)

// AuditLogLevel 操作日志级别
type AuditLogLevel int

const (
	MinorLevel     AuditLogLevel = 0 // 操作日志中操作类的次要级别
	ImportantLevel AuditLogLevel = 1 // 操作日志中操作类的重要级别
	LogLevelAuto   AuditLogLevel = 3 // 操作日志中查询类的自动查询
	LogLevelManual AuditLogLevel = 4 // 操作日志中查询类的手动查询
)

type OperateType int

const (
	GET OperateType = iota
	ADD
	MOD
	DELETE
	DOWNLOAD
	UPLOAD
	UPHOLD
)

// AuditsInfo 审计日志信息
type AuditsInfo struct {
	Terminal string
	UserName string
	Detail   string
	DetailZh string
}

// AuditsPara 鉴权函数参数
type AuditsPara struct {
	OperationZH string
	OperationEN string
	OperateType OperateType
	Level       AuditLogLevel
	Username    string
	Terminal    string
	Result      int
	Detail      string
	DetailZH    string
}

// AuditsLog  审计日志函数
func AuditsLog(auditsPara *AuditsPara, requestURL string) {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	body := make(map[string]interface{})

	operation := make(map[string]string)
	operation["OP_ZH"] = auditsPara.OperationZH
	operation["OP_EN"] = auditsPara.OperationEN
	opstr, err := json.Marshal(operation)
	if err != nil {
		return
	}
	body["operation"] = string(opstr)
	body["level"] = auditsPara.Level
	if requestURL == SecLog {
		body["level"] = ImportantLevel
	}
	body["userName"] = auditsPara.Username
	body["dateTime"] = time.Now().UnixNano() / 1e6
	body["appName"] = os.Getenv("APPNAME")
	body["appId"] = os.Getenv("APPID")
	body["terminal"] = auditsPara.Terminal
	body["serviceName"] = constants.ServiceName
	body["result"] = auditsPara.Result
	body["detail"] = auditsPara.Detail
	body["detail_zh"] = auditsPara.DetailZH
	if requestURL == OpsLog {
		body["operateType"] = strconv.Itoa(int(auditsPara.OperateType))
	}

	// Auditlog 服务端需要接受字符串序列化成的json，所以序列化两次
	bs, err := json.Marshal(body)
	if err != nil {
		Errorf("marshal audit log body failed, err: %v", err)
		return
	}
	bs2, err := json.Marshal(string(bs))
	if err != nil {
		Errorf("marshal json string of audit log body failed, err: %v", err)
		return
	}

	Infof("upload audit log %s", string(bs2))

	resp, err := gsfapi.NewCspRestInvoker().Invoke(http.MethodPost, requestURL, headers, bs2)
	if resp != nil {
		defer resp.Close()
	}

	if resp == nil || err != nil {
		Errorf("send auditlog faild, Cse RestInvoker failed, *rest.Response is nil, err: %v", err)
		return
	}

	statusCode := resp.GetStatusCode()
	if statusCode != http.StatusOK {
		Errorf("send auditlog faild, responseBody is %s, err: %v", string(resp.ReadBody()), err)
	}
	Infof("upload audit log success")
}

// AuditsSecAndOpsLog 同时记录安全日志和操作日志
func AuditsSecAndOpsLog(secAuditsPara, opsAuditsPara *AuditsPara) {
	AuditsLog(secAuditsPara, SecLog)
	AuditsLog(opsAuditsPara, OpsLog)
}
