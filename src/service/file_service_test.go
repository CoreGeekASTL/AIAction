// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package service
package service

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

/*
* 测试用例描述： TestCleanFileName
* 预置条件：无
* 操作步骤：
*     1. 测试正常文件名
*     2. 测试空文件名
*     3. 测试包含斜杠的文件名
*     4. 测试包含反斜杠的文件名
*     5. 测试包含单个斜杠（根目录路径）的文件名
* 预期结果：
*     1. 返回错误为 nil
*     2. 返回错误为 nil
*     3. 返回错误
*     4. 返回错误
*     5. 返回错误
* 修改历史：
*     1. 2025-8-27 新建测试用例
 */
func TestCleanFileName(t *testing.T) {
	fileService := NewFileService()
	tests := []struct {
		name       string
		fileName   string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "test normal file name",
			fileName:   "test.txt",
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "test empty file name",
			fileName:   "",
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "test file name with slash",
			fileName:   "test/evil/path.txt",
			wantErr:    true,
			wantErrMsg: "invalid file name",
		},
		{
			name:       "test file name with backslash",
			fileName:   "test\\evil\\path.txt",
			wantErr:    true,
			wantErrMsg: "invalid file name",
		},
		{
			name:       "test file name with root path",
			fileName:   "/",
			wantErr:    true,
			wantErrMsg: "invalid file name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fileService.cleanFileName(tt.fileName)
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
* 测试用例描述： TestUploadFile
* 预置条件：无
* 操作步骤：
*     1. 测试文件名包含路径分隔符
* 预期结果：
*     1. 返回错误
* 修改历史：
*     1. 2025-8-27 新建测试用例
*     2. 2025-11-07 修复测试用例
 */
func TestUploadFile(t *testing.T) {
	fileService := NewFileService()
	tests := []struct {
		bucket      string
		name        string
		fileName    string
		fileContent []byte
		wantErr     bool
		wantPath    string
	}{
		{
			bucket:      "file",
			name:        "test upload file with path separator",
			fileName:    "test/evil/path.txt",
			fileContent: []byte("test content"),
			wantErr:     true,
			wantPath:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, err := fileService.UploadFile(tt.bucket, tt.fileName, tt.fileContent)
			assert.Equal(t, tt.wantPath, filePath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

/*
* 测试用例描述： TestDownloadFile
* 预置条件：
* 操作步骤：
*     1. 测试文件名包含路径分隔符
* 预期结果：
*     1. 返回错误
* 修改历史：
*     1. 2025-8-27 新建测试用例
*     2. 2025-11-07 修复测试用例
 */
func TestDownloadFile(t *testing.T) {
	fileService := NewFileService()
	tests := []struct {
		bucket   string
		name     string
		fileName string
		wantErr  bool
	}{
		{
			bucket:   "file",
			name:     "test download file with path separator",
			fileName: "test/evil/path.txt",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := fileService.DownloadFile(tt.bucket, tt.fileName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, file)
				return
			}
			// 检查错误
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}
