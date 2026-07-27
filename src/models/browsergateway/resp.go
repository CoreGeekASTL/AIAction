// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package browsergateway
package browsergateway

type ExtensionLoadResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    ExtensionLoadRequest `json:"data"`
}
