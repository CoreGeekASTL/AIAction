// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package req
package req

import (
	"fmt"
	"mime/multipart"

	"GIDS/common/constants"
)

type UploadPluginPackageReq struct {
	Filename string `json:"filename" json:"filename,omitempty"`
	File     multipart.File
	Size     int64
}

func (r *UploadPluginPackageReq) Validate() error {
	if r.Filename == "" || r.File == nil || r.Size == 0 || r.Size > constants.MaxFileSize {
		return fmt.Errorf("invalid param")
	} else {
		return nil
	}
}

type PluginPackageReq struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

func (r *PluginPackageReq) GetKey() string {
	return fmt.Sprintf("%s:%s:%s", r.Type, r.Name, r.Version)
}

func (r *PluginPackageReq) Validate() error {
	if r.Name == "" || r.Version == "" || r.Type == "" {
		return fmt.Errorf("invalid param")
	}
	return nil
}
