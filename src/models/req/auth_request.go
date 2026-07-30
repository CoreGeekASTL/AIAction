// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package req
package req

import "errors"

// AuthIMEIRequest 终端联合鉴权请求，BGW 反调 /auth/v1/authIMEI 使用
type AuthIMEIRequest struct {
	IMEI string `json:"imei"`
	IMSI string `json:"imsi"`
}

func (r *AuthIMEIRequest) Validate() error {
	if r.IMEI == "" || r.IMSI == "" {
		return errors.New("imei or imsi cannot be empty")
	}
	return nil
}
