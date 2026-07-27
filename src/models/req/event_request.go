// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

package req

type ClientEventRequest struct {
	HSMan   string `json:"hsman"`
	HSType  string `json:"hstype"`
	AppType string `json:"appType"`
	IMEI    string `json:"imei"`
	IMSI    string `json:"imsi"`
	Type    string `json:"type"`
}

func (e ClientEventRequest) Validate() error {
	return nil
}

type AppUseTimesEvent struct {
	UseTimes string `json:"useTimes"`
	HSMan    string `json:"hsman"`
	HSType   string `json:"hstype"`
	EXTType  string `json:"exttype"`
	AppType  string `json:"appType"`
	AppId    string `json:"appId"`
	SCHeight string `json:"scheight"`
	SCWidth  string `json:"scwidth"`
	IMEI     string `json:"imei"`
	IMSI     string `json:"imsi"`
	PlayMode string `json:"playMode"`
}

func (e AppUseTimesEvent) Validate() error {
	return nil
}
