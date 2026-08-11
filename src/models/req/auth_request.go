// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package req

import "errors"

// AuthIMEIRequest 终端联合鉴权请求
type AuthIMEIRequest struct {
	IMEI string `json:"imei"`
	IMSI string `json:"imsi"`
}

// Validate 仅做非空检查，15 位格式校验归 Service 判定
func (r AuthIMEIRequest) Validate() error {
	if r.IMEI == "" || r.IMSI == "" {
		return errors.New("imei and imsi cannot be empty")
	}
	return nil
}
