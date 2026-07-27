// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package resp
package resp

type DataResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

type GridLoginAuthResponse struct {
	BaseResponse
	Data struct {
		Token              string `json:"token"`
		ExpiresTime        string `json:"expiresTime"`
		NodeGateWayURL     string `json:"nodeGateWayUrl"`
		NodeIntranetWayURL string `json:"nodeIntranetWayUrl"`
		NodeCapacity       int32  `json:"nodeCapacity"`
		TimeAxis           int64  `json:"timeAxis"`
	} `json:"data"`
}

type DeviceLoginAuthResponse struct {
	BaseResponse
	Data LoginInfo `json:"data"`
}

type AuthInfo struct {
	Token       string `json:"token"`
	ExpiresTime int64  `json:"expiresTime"`
	TimeAxis    int64  `json:"timeAxis"`
}

type AssignInfo struct {
	TcpAddr             string `json:"tcpAddr,omitempty"`
	TlsTcpAddr          string `json:"tlsTcpAddr,omitempty"`
	VideoMode           int    `json:"videoMode,omitempty"`
	ShortAddr           string `json:"shortAddr,omitempty"`
	NodeGateWayURL      string `json:"nodeGateWayUrl"`
	HttpsShortAddr      string `json:"httpsShortAddr"`
	HttpsNodeGateWayUrl string `json:"httpsNodeGateWayUrl"`
	NodeIntranetWayURL  string `json:"nodeIntranetWayUrl"`
	NodeCapacity        int    `json:"nodeCapacity,omitempty"`
}

type LoginInfo struct {
	AuthInfo
	AssignInfo
}
