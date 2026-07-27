// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
package oss

import (
	"context"
	"testing"

	"GIDS/common/conf"
)

/*
* 测试用例描述：测试 GetObject 方法是否能正确获取对象
* 预置条件：
*		1. 服务正常运行
* 操作步骤：
* 		1. 测试有效参数
* 		2. 测试无效参数 - 空的 bucketName
* 		3. 测试无效参数 - bucketName 包含无效字符
* 		4. 测试无效参数 - bucketName 少于 3 个字符
* 		5. 测试无效参数 - bucketName 多于 63 个字符
* 		6. 测试无效参数 - bucketName 为 IP 地址
* 		7. 测试无效参数 - 空的 objectName
* 		7. 测试无效参数 - objectName 多于 1024 个字符
* 		8. 测试有效参数 - objectName 包含中文字符串
* 		9. 测试无效参数 - objectName 包含非 UTF-8 字符串
* 预期结果：
*       1. 返回错误为 nil
*       2. 返回错误，错误信息为"Bucket name cannot be empty"
*       3. 返回错误，错误信息为"Bucket name contains invalid characters"
*       4. 返回错误，错误信息为"Bucket name cannot be shorter than 3 characters"
*       5. 返回错误，错误信息为"Bucket name cannot be longer than 63 characters"
*       6. 返回错误，错误信息为"Bucket name cannot be an ip address"
*       7. 返回错误，错误信息为"Object name cannot be empty"
*       8. 返回错误为 nil
*       9. 返回错误，错误信息为"Object name cannot be empty"
* 修改历史：
*       1. 2025-8-26 新建测试用例
 */
func TestGetObject(t *testing.T) {
	// 从配置中获取Endpoint
	config := conf.Instance().OSS.Endpoint
	// 初始化Client
	if err := Init(conf.OSSConfig{
		Endpoint: config,
	}); err != nil {
		t.Fatalf("failed to initialize client: %v", err)
	}

	type args struct {
		ctx     context.Context
		options GetObjectOptions
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success case - valid request",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "upload-bucket",
					FileName:   "test.txt",
				},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "error case - empty bucket name",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "",
					FileName:   "test.txt",
				},
			},
			wantErr:    true,
			wantErrMsg: "Bucket name cannot be empty",
		},
		{
			name: "error case - bucket name contains invalid characters",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "/upload-bucket",
					FileName:   "test.txt",
				},
			},
			wantErr:    true,
			wantErrMsg: "Bucket name contains invalid characters",
		},
		{
			name: "error case - bucket name shorter than 3 characters",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "11",
					FileName:   "test.txt",
				},
			},
			wantErr:    true,
			wantErrMsg: "Bucket name cannot be shorter than 3 characters",
		},
		{
			name: "error case - bucket name longer than 63 characters",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "1111111111111111111111111111111111111111111111111111111111111111",
					FileName:   "test.txt",
				},
			},
			wantErr:    true,
			wantErrMsg: "Bucket name cannot be longer than 63 characters",
		},
		{
			name: "error case - bucket name is an ip address",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "127.0.0.1",
					FileName:   "test.txt",
				},
			},
			wantErr:    true,
			wantErrMsg: "Bucket name cannot be an ip address",
		},
		{
			name: "error case - empty object name",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "upload-bucket",
					FileName:   "",
				},
			},
			wantErr:    true,
			wantErrMsg: "Object name cannot be empty",
		},
		{
			name: "error case - Object name longer than 1024 characters",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "upload-bucket",
					FileName:   "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
				},
			},
			wantErr:    true,
			wantErrMsg: "Object name cannot be longer than 1024 characters",
		},
		{
			name: "error case - Object name has Chinese strings",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "upload-bucket",
					FileName:   "Hello, 世界",
				},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "error case - Object name has non UTF-8 strings",
			args: args{
				ctx: context.Background(),
				options: GetObjectOptions{
					BucketName: "upload-bucket",
					FileName:   string([]byte{0xff, 0xfe, 0xfd}),
				},
			},
			wantErr:    true,
			wantErrMsg: "Object name with non UTF-8 strings are not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			object, err := Instance().GetObject(tt.args.ctx, tt.args.options)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("expected error %q, got %q", tt.wantErrMsg, err.Error())
				}
				// 确保关闭返回的 object ，避免资源泄漏
				if object != nil {
					if closeErr := object.Close(); closeErr != nil {
						t.Errorf("failed to close object: %v", closeErr)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if object == nil {
					t.Error("expected non-nil ReadCloser, got nil")
				} else {
					defer object.Close() // 测试结束后确保关闭资源，避免引发文件句柄泄漏
				}
			}
		})
	}
}
