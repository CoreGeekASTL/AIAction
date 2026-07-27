/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package https
package https

import (
	"net/http"
	"time"

	"Go-chassis-extend/api/GSF/sdk/tlsutil"
	"Go-chassis-extend/api/ServiceComb/go-chassis/core/common"

	"GIDS/common/logger"
)

var client HTTPDoer
var muenClient HTTPDoer
var innerClient HTTPDoer

const timeout = 120
const idleTimeout = 300
const timeoutMultiplier = 2

func Init() {
	client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: timeout * time.Second,
			TLSHandshakeTimeout:   timeout * time.Second,
			ExpectContinueTimeout: timeout * time.Second,
			IdleConnTimeout:       idleTimeout * time.Second,
		},
		Timeout: timeout * timeoutMultiplier * time.Second,
	}
	initInnerClient()
}

func Instance() HTTPDoer {
	return client
}

// InnerInstance ge inner client
func InnerInstance() HTTPDoer {
	return innerClient
}

// 初始化内部通信clent，需要在csp框架初始化完成之后
func initInnerClient() {
	tlsConfig, err := tlsutil.GetTLSConfig("registry", "", common.Consumer)
	if err != nil {
		logger.Errorf("failed to get inner client tlsConfig:%v", err)
		return
	}
	innerClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}

func MuenInstance() HTTPDoer {
	return muenClient
}

// InitMuenClient 对外通信 使用沐恩提供的ca证书(外部证书)
func InitMuenClient() {
	muenClient = newHttpsClient(CertInfo{})
	logger.Infof("init client for muen server end")
}

func MuenCertUpdate(info CertInfo) {
	muenClient = newHttpsClient(info)
}

func newHttpsClient(info CertInfo) HTTPDoer {
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       GetTLS(info, ClientType),
			ResponseHeaderTimeout: timeout * time.Second,
			TLSHandshakeTimeout:   timeout * time.Second,
			ExpectContinueTimeout: timeout * time.Second,
			IdleConnTimeout:       idleTimeout * time.Second,
		},
		Timeout: timeout * timeoutMultiplier * time.Second,
	}
	return cli
}
