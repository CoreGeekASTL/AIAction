// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

package browsergateway

import "GIDS/models/db"

const percentageMultiplier = 100

type ServiceInstance struct {
	BrowserInnerEndpoint     string          `json:"browserGWInnerEndpoint,omitempty"`
	MediaExtendEndpoint      string          `json:"edgeMediaExtendEndpoint,omitempty"`
	MediaInnerEndpoint       string          `json:"edgeMediaInnerEndpoint,omitempty"`
	ControlExtendEndpoint    string          `json:"edgeControlExtendEndpoint,omitempty"`
	MediaTlsExtendEndpoint   string          `json:"edgeMediaTlsExtendEndpoint,omitempty"`
	ControlTlsExtendEndpoint string          `json:"edgeControlTlsExtendEndpoint,omitempty"`
	Cap                      int             `json:"cap"`
	Used                     int             `json:"used"`
	PluginStatus             db.ActiveStatus `json:"pluginStatus"`
	IsHealthy                bool            `json:"isHealthy,omitempty"`
	CheckMsg                 string          `json:"checkMsg,omitempty"`
}

func (s *ServiceInstance) GetKey() string {
	return s.BrowserInnerEndpoint
}

type ServiceInstanceList []ServiceInstance

func (s ServiceInstanceList) Len() int {
	return len(s)
}

func (s ServiceInstanceList) Less(i, j int) bool {
	iPercent := s[i].Used * percentageMultiplier / s[i].Cap
	jPercent := s[j].Used * percentageMultiplier / s[j].Cap
	return iPercent < jPercent
}

func (s ServiceInstanceList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
