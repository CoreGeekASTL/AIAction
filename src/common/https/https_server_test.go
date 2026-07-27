// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package https pkg about https
package https

import (
	"fmt"
	"os"
	"testing"

	"code.huawei.com/paaslite/dev-tool/gomockit"
	"github.com/stretchr/testify/assert"
)

/*
 * 测试用例描述：Test GetLocalIP
 * 预制条件：
 *         1. 无
 * 操作步骤：
 *         1. 调用 GetLocalIP
 * 预期结果：
 *         1. take precedence use first parameter
 *         2. use default eth as result when first parameter is empty
 *         3. failed when both env is empty
 * 修改历史：
 *         1. 2025/12/16 新建测试用例
 */
func TestGetLocalIP(t *testing.T) {
	t.Skip("skip: Linux network interface not available on Windows")
	gomockit.MockFunc(getEthIP, func(name string) (string, error) {
		if name == "bond-fabric" {
			return "192.168.1.1", nil
		}
		if name == "bond-base" {
			return "192.168.2.1", nil
		}
		return "", fmt.Errorf("invalid eth")
	})
	defer gomockit.Reset()
	type args struct {
		ethName    string
		defaultEth string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "take precedence use first parameter",
			args:    args{ethName: "bond-fabric", defaultEth: "bond-base"},
			want:    "192.168.1.1",
			wantErr: assert.NoError,
		},
		{
			name:    "use default eth as result when first parameter is empty",
			args:    args{ethName: "", defaultEth: "bond-base"},
			want:    "192.168.2.1",
			wantErr: assert.NoError,
		},
		{
			name:    "failed when both env is empty",
			args:    args{ethName: "", defaultEth: ""},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ethName != "" {
				if err := os.Setenv(tt.args.ethName, tt.args.ethName); err != nil {
					t.Logf("setenv failed: %v", err)
				}
				defer os.Unsetenv(tt.args.ethName)
			}
			if tt.args.defaultEth != "" {
				if err := os.Setenv(tt.args.defaultEth, tt.args.defaultEth); err != nil {
					t.Logf("setenv failed: %v", err)
				}
				defer os.Unsetenv(tt.args.defaultEth)
			}
			got, err := GetLocalIP(tt.args.ethName, tt.args.defaultEth)
			if !tt.wantErr(t, err, fmt.Sprintf("GetLocalIP(%v, %v)", tt.args.ethName, tt.args.defaultEth)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetLocalIP(%v, %v)", tt.args.ethName, tt.args.defaultEth)
		})
	}

}
