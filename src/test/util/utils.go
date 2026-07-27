// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package util 提供测试辅助工具函数。
package util

// It 补充goconvey写UT时确实的串行流程， 用于区分串行流程操作
func It(text string, f func()) {
	f()
}
