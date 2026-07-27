// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package service
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"GIDS/common/logger"
)

const defaultCacheTimeoutSeconds = 5

// DeleteCache 删除页面缓存处理服务
var DeleteCache = DeleteCacheImpl

func DeleteCacheImpl(imei, imsi string) error {
	if imei == "" || imsi == "" {
		return errors.New("IMEI or IMSI cannot be empty")
	}

	// 获取 BrowserGW 实例列表
	browserService := NewBrowserService()
	instances := browserService.GetAllReadyServiceInstances()

	if len(instances) == 0 {
		return errors.New("no BrowserGW instances available")
	}

	for _, instance := range instances {
		// 调用 BrowserGW 接口删除对象存储缓存
		err := callBrowserGW(instance.BrowserInnerEndpoint, imei, imsi)
		if err != nil {
			logger.Errorf("failed to call BrowserGW interface: %v", err)
		}
	}

	return nil
}

// callBrowserGW 调用 BrowserGW 接口
func callBrowserGW(browserGW string, imei, imsi string) error {
	client := &http.Client{
		Timeout: time.Second * defaultCacheTimeoutSeconds,
	}

	// 构建请求体
	requestBody := struct {
		Imei string `json:"imei"`
		Imsi string `json:"imsi"`
	}{
		Imei: imei,
		Imsi: imsi,
	}

	// 将请求体序列化为 JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 构建请求 URL
	reqURL := fmt.Sprintf("http://%s/browsergw/browser/userdata/delete", browserGW)

	req, err := http.NewRequest("DELETE", reqURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送 DELETE 请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send DELETE request to BrowserGW: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("browserGW returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
