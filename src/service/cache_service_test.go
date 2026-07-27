// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package service
package service

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
* 测试用例描述：TestDeleteCache
* 预置条件：
* 操作步骤：
*     1. 测试参数为空
* 预期结果：
*     1. 返回错误，参数为空
* 修改历史：
*     1. 2025-8-27 新建测试用例
 */
func TestDeleteCache(t *testing.T) {
	tests := []struct {
		name       string
		imei       string
		imsi       string
		redisErr   error
		browserErr error
		wantErr    bool
	}{
		{
			name:       "test delete cache empty params",
			imei:       "",
			imsi:       "",
			redisErr:   nil,
			browserErr: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用被测函数
			err := DeleteCacheImpl(tt.imei, tt.imsi)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

/*
* 测试用例描述：TestCallBrowserGW
* 预置条件：
* 操作步骤：
*     1. 测试调用 BrowserGW 失败
* 预期结果：
*     1. 返回 500 错误码，调用失败
* 修改历史：
*     1. 2025-8-27 新建测试用例
 */
func TestCallBrowserGW(t *testing.T) {
	tests := []struct {
		name       string
		browserGW  string
		imei       string
		imsi       string
		respStatus int
		respBody   string
		wantErr    bool
	}{
		{
			name:       "test call browsergw failed",
			browserGW:  "http://fake-browsergw:8080",
			imei:       "1234567890",
			imsi:       "9876543210",
			respStatus: http.StatusInternalServerError,
			respBody:   "internal server error",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建fake HTTP client
			client := &http.Client{}
			// 创建fake response
			resp := &http.Response{
				StatusCode: tt.respStatus,
				Body:       io.NopCloser(bytes.NewBufferString(tt.respBody)),
			}
			// 创建fake roundtripper
			fakeTransport := &fakeRoundTripper{
				resp: resp,
				err:  nil,
			}
			client.Transport = fakeTransport

			// 调用被测函数
			err := callBrowserGW(tt.browserGW, tt.imei, tt.imsi)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// fakeRoundTripper 用于测试的fake HTTP传输实现
type fakeRoundTripper struct {
	resp *http.Response
	err  error
}

func (f *fakeRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return f.resp, f.err
}
