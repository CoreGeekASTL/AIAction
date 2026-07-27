// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
package req

import (
	"testing"
)

/*
* 测试用例描述：测试 LoginAuthRequest 、 SyncBrowserConfigRequest 和 DeleteCacheRequest 的 Validate 方法
* 预置条件：
*		无
* 操作步骤：
* 		1. 调用 Validate 方法
* 预期结果：
*       1. Validate 方法返回 nil
* 修改历史：
*       1. 2025-8-26 新建测试用例
 */
func TestLoginAuthRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request LoginAuthRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			request: LoginAuthRequest{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*
* 测试用例描述：测试 SyncBrowserConfigRequest 的 Validate 方法是否能正确验证请求参数
* 预置条件：
*		1. 服务正常运行
* 操作步骤：
* 		1. 发送有效的 SyncBrowserConfigRequest 请求
* 预期结果：
*       1. 成功通过验证，返回 nil
* 修改历史：
*       1. 2025-8-26 新建测试用例
 */
func TestSyncBrowserConfigRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request SyncBrowserConfigRequest
		wantErr bool
	}{
		{
			name:    "success case - valid request",
			request: SyncBrowserConfigRequest{RouteAPPConfigList: []RouteAppConfig{{Name: "test"}}, ChromeConfigList: []ChromeConfig{{Manufacturer: "test"}}, URLConfigs: []URLConfigs{{URL: "test"}}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncBrowserConfigRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*
* 测试用例描述：测试 DeleteCacheRequest 的 Validate 方法
* 预置条件：
*		1. 服务正常运行
* 操作步骤：
* 		1. 测试有效参数
* 		2. 测试无效参数（缺少 IMEI ）
* 		3. 测试无效参数（缺少 IMSI ）
* 		4. 测试无效参数（ IMEI 和 IMSI 都为空）
* 预期结果：
*       1. 返回错误为 nil
*       2. 返回错误，错误信息为"IMEI cannot be empty"
*       3. 返回错误，错误信息为"IMSI cannot be empty"
*       4. 返回多个错误
* 修改历史：
*       1. 2025-8-26 新建测试用例
 */
func TestDeleteCacheRequestValidate(t *testing.T) {
	tests := []struct {
		name       string
		request    DeleteCacheRequest
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "success case - valid request",
			request:    DeleteCacheRequest{IMEI: "123456789012345", IMSI: "123456789012345"},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "error case - empty imei",
			request:    DeleteCacheRequest{IMSI: "123456789012345"},
			wantErr:    true,
			wantErrMsg: "IMEI cannot be empty",
		},
		{
			name:       "error case - empty imsi",
			request:    DeleteCacheRequest{IMEI: "123456789012345"},
			wantErr:    true,
			wantErrMsg: "IMSI cannot be empty",
		},
		{
			name:       "error case - both empty",
			request:    DeleteCacheRequest{},
			wantErr:    true,
			wantErrMsg: "IMEI cannot be empty; IMSI cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("expected error %q, got %q", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

/*
* 测试用例描述：测试 FileUploadRequest 的 Validate 方法
* 预置条件：
*		1. 服务正常运行
* 操作步骤：
* 		1. 测试有效参数
* 		2. 测试无效参数（缺少 fileName ）
* 预期结果：
*       1. 返回 nil
*       2. 返回错误，错误信息为"fileName cannot be empty"
* 修改历史：
*       1. 2025-8-26 新建测试用例
 */
func TestFileUploadRequestValidate(t *testing.T) {
	tests := []struct {
		name       string
		request    FileUploadRequest
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "success case - valid request",
			request:    FileUploadRequest{FileName: "test.txt"},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "error case - empty fileName",
			request:    FileUploadRequest{FileName: ""},
			wantErr:    true,
			wantErrMsg: "fileName cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("expected error %q, got %q", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
