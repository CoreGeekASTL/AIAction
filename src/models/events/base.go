// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package events
package events

import (
	"encoding/json"
	"time"

	"GIDS/common/constants"
	"GIDS/common/logger"
)

type EventType string

type EventDesc struct {
	Event        EventType `json:"event"`
	EventDesc    string    `json:"eventDesc"`
	EventTrigger string    `json:"eventTrigger"`
}

const (
	Login       EventType = "browser_user_http_login"
	LoginError  EventType = "browser_user_http_login_error"
	Client      EventType = "browser_client_error"
	AppUseTimes EventType = "browser_client_app_use_times"
)

const (
	CloudBrowser = "cloud-browser"
)

var (
	eventTypeMap = map[EventType]EventDesc{
		Login: {
			Event:        "browser_user_http_login",
			EventDesc:    "云浏览器用户HTTP登录埋点",
			EventTrigger: "client",
		},
		Client: {
			Event:        "browser_client_error",
			EventDesc:    "客户端异常",
			EventTrigger: "client",
		},
		AppUseTimes: {
			Event:        "browser_client_app_use_times",
			EventDesc:    "客户端上传应用的使用时长",
			EventTrigger: "client",
		},
	}
)

type Info struct {
	EventDesc
	Service   string      `json:"service"`
	EventTime string      `json:"eventTime"`
	Env       string      `json:"env"`
	Hostname  string      `json:"hostname"`
	Object    string      `json:"object"`
	EventData interface{} `json:"eventData"`
}

func NewInfo(event EventType) *Info {
	return &Info{
		EventDesc: eventTypeMap[event],
		Service:   constants.ComponentName,
		EventTime: time.Now().Format("2006-01-02 15:04:05"),
		Env:       "",
		Hostname:  "",
		Object:    CloudBrowser,
		EventData: nil,
	}
}

func (i *Info) SetEventData(eventData interface{}) {
	i.EventData = eventData
}

// ToJSON func
func (i *Info) ToJSON() []byte {
	content, err := json.Marshal(i)
	if err != nil {
		logger.Errorf("convert to json failed, source: [%+v]", i)
		return []byte("")
	}

	return content
}

type LoginEventData struct {
	IMEI                string `json:"imei"`
	IMSI                string `json:"imsi"`
	AppType             string `json:"appType"`
	ExtType             string `json:"extType"`
	HSMan               string `json:"hsman"`
	HSType              string `json:"hstype"`
	LoginTime           string `json:"loginTime"`
	TcpAddr             string `json:"tcpAddr"`
	TlsTcpAddr          string `json:"tlsTcpAddr"`
	VideoMode           int    `json:"videoMode"`
	ShortAddr           string `json:"shortAddr"`
	NodeGateWayUrl      string `json:"nodeGateWayUrl"`
	HttpsShortAddr      string `json:"httpsShortAddr"`
	HttpsNodeGateWayUrl string `json:"httpsNodeGateWayUrl"`
	TotalKb             string `json:"totalKb"`
	FreeKb              string `json:"freeKb"`
}

type ClientEventData struct {
	HSMan   string `json:"hsman"`
	HSType  string `json:"hstype"`
	AppType string `json:"appType"`
	IMEI    string `json:"imei"`
	IMSI    string `json:"imsi"`
	Type    string `json:"type"`
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
