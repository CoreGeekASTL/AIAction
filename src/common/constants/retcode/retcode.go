// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package retcode 定义HTTP返回码和认证状态码常量。
package retcode

const (
	Success        = 200
	InternalFailed = -1
	ClientFailed   = -2
	AuthPassed     = 200
	AuthFailed     = 401
)
