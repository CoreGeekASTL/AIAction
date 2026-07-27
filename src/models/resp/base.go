// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package resp
package resp

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}
