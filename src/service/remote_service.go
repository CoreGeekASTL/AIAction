/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package service

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/https"
	"GIDS/common/logger"
	"GIDS/models/req"
	"GIDS/models/resp"
)

func MuenDeviceLogin(request *req.LoginAuthRequest) *resp.LoginInfo {
	muenConfigCenter := NewConfigCenterService()
	enableHttps, endpoint := getMoonConfigUrl(muenConfigCenter)
	client := https.Instance()
	url := fmt.Sprintf("%s/app-api/devicetcp/app/login/v1/deviceLoginAuth", endpoint)
	if enableHttps {
		url = "https://" + url
		client = https.MuenInstance()
	} else {
		url = "http://" + url
	}
	logger.Infof("use cloud service url %v", url)
	response := https.NewRequest(client).WithRetry(defaultRetryCount).
		Method("POST").
		URL(url).
		ParamFromInterface(request).Complete().Do()
	if response.Error() != nil {
		logger.Errorf("login in muen failed: %v", response.Error())
		return nil
	}
	if !response.IsSuccessCode() || response.Error() != nil {
		logger.Errorf("login in muen failed, status is %d, err is %v",
			response.StatusCode(), response.Error())
		return nil
	}
	authResponse := &resp.DeviceLoginAuthResponse{}
	err := response.ResponseToStruct(authResponse)
	logger.Infof("result from muen : %v", authResponse.Data)
	if err != nil {
		logger.Errorf("login in muen failed: %v", err)
		return nil
	}
	logger.Infof("respone from muen: %v", authResponse)
	return &authResponse.Data
}

func getMoonConfigUrl(configCenter ConfigCenterService) (bool, string) {
	cfgUrl := beego.AppConfig.DefaultString("moon::titokEndpoint", "")
	config, ok := configCenter.GetConfig("moon::titokEndpoint")
	if ok && config != "" {
		cfgUrl = config
	}

	cfgHttpsUrl := beego.AppConfig.DefaultString("moon::httpsTitokEndpoint", "")
	httpsConfig, ok := configCenter.GetConfig("moon::httpsTitokEndpoint")
	if ok && httpsConfig != "" {
		cfgHttpsUrl = httpsConfig
	}

	enableHttps := beego.AppConfig.DefaultString("moon::enableHttps", "")
	enableHttpsConfig, ok := configCenter.GetConfig("moon::enableHttps")
	if ok && enableHttpsConfig != "" {
		enableHttps = enableHttpsConfig
	}

	if enableHttps == "true" && cfgHttpsUrl != "" {
		return true, cfgHttpsUrl
	}
	return false, cfgUrl
}
