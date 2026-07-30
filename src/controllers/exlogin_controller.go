// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package controllers include web controllers
package controllers

import (
	"time"

	"GIDS/common/constants"
	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/models/events"
	"GIDS/models/req"
	"GIDS/models/resp"
	"GIDS/service"
)

// ExLoginController for external communication, use https
type ExLoginController struct {
	BaseController
	userService    service.UserService
	browserService service.BrowserService
	eventService   service.EventService
	authService    service.AuthService
}

func (c *ExLoginController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/app-api/devicetcp/app/login/v1/gridLoginAuth":            "POST:GridLoginAuth",
			"/app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser": "POST:GridLoginAuthOpenBrowser",
			"/app-api/devicetcp/app/login/v1/deviceLoginAuth":          "POST:DeviceLoginAuth",
		},
	}
}

func (c *ExLoginController) Prepare() {
	c.userService = service.NewUserService()
	c.browserService = service.NewBrowserService()
	c.eventService = service.NewEventService()
	c.authService = service.NewAuthService()
}

func (c *ExLoginController) GridLoginAuth() {
	_, response := c.loginAuth(false)
	if response == nil {
		return
	}
	response.Data.TcpAddr = ""
	response.Data.TlsTcpAddr = ""
	response.Data.VideoMode = 0
	response.Data.ShortAddr = ""
	response.Data.HttpsShortAddr = ""
	response.Data.NodeIntranetWayURL = ""
	c.OK(response)
}

func (c *ExLoginController) GridLoginAuthOpenBrowser() {
	_, response := c.loginAuth(true)
	if response == nil {
		return
	}
	response.Data.TcpAddr = ""
	response.Data.TlsTcpAddr = ""
	response.Data.VideoMode = 0
	response.Data.ShortAddr = ""
	response.Data.HttpsShortAddr = ""
	response.Data.NodeIntranetWayURL = ""
	c.OK(response)
}

func (c *ExLoginController) DeviceLoginAuth() {
	request, response := c.loginAuth(false)
	if response == nil {
		return
	}
	response.Data.NodeIntranetWayURL = ""
	// tiktok登录 取muen token
	if request.AppType == constants.TikTokAppType {
		muenLogin := service.MuenDeviceLogin(request)
		if muenLogin == nil {
			logger.Errorf("login in muen failed")
			c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: "login failed"})
			return
		}
		if err := c.browserService.UpdateUserToken(muenLogin, request); err != nil {
			logger.Errorf("update user token failed")
			c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: "login failed"})
			return
		}
		response.Data = *muenLogin
	}
	c.OK(response)
}

func (c *ExLoginController) loginAuth(preOpenBrowser bool) (*req.LoginAuthRequest, *resp.DeviceLoginAuthResponse) {
	request := new(req.LoginAuthRequest)
	err := c.RequestBodyUnmarshalTo(request)
	if err != nil {
		logger.Infof("[loginAuth] unmarshal failed, err: [%v], request: [%v]", err, request)
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return nil, nil
	}

	if pass, _ := c.authService.Check(request.IMEI, request.IMSI); !pass {
		logger.Warnf("[loginAuth] terminal auth rejected")
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: "auth rejected"})
		return nil, nil
	}

	err = c.userService.CreateOrUpdateUser(request)
	if err != nil {
		logger.Errorf("[loginAuth] create or update user failed, err is [%v]", err)
		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: "Login failed"})
		return nil, nil
	}
	LoginInfo, err := c.browserService.RouteToInstance(request)
	if err != nil {
		logger.Warnf("[loginAuth] route to instance failed, err is [%v], auth passed, return empty instance", err)
		LoginInfo = resp.LoginInfo{}
	}
	if preOpenBrowser {
		c.browserService.PreOpenBrowser(request)
	}
	response := resp.DeviceLoginAuthResponse{
		BaseResponse: resp.BaseResponse{
			Code:    retcode.Success,
			Message: "success",
		},
		Data: LoginInfo,
	}
	c.reportDeviceLoginEvent(request, response)
	return request, &response
}

func (c *ExLoginController) reportDeviceLoginEvent(request *req.LoginAuthRequest, response resp.DeviceLoginAuthResponse) {
	event := events.NewInfo(events.Login)
	event.SetEventData(events.LoginEventData{
		IMEI:                request.IMEI,
		IMSI:                request.IMSI,
		AppType:             request.AppType,
		ExtType:             request.ExtendModel,
		HSMan:               request.Manufacturer,
		HSType:              request.Model,
		LoginTime:           time.Unix(response.Data.TimeAxis, 0).Format("2006-01-02 15:04:05"),
		TcpAddr:             response.Data.TcpAddr,
		TlsTcpAddr:          response.Data.TlsTcpAddr,
		VideoMode:           response.Data.VideoMode,
		ShortAddr:           response.Data.ShortAddr,
		NodeGateWayUrl:      response.Data.NodeGateWayURL,
		HttpsShortAddr:      response.Data.HttpsShortAddr,
		HttpsNodeGateWayUrl: response.Data.HttpsNodeGateWayUrl,
		TotalKb:             request.TotalKb,
		FreeKb:              request.FreeKb,
	})

	logger.Infof("get TlsTcpAddr %v", response.Data.TlsTcpAddr)
	logger.Infof("get HttpsNodeGateWayUrl %v", response.Data.HttpsNodeGateWayUrl)

	// #toDo Event记录失败导致的异常暂不导致Login失败
	if err := c.eventService.ReportEvent(event); err != nil {
		logger.Errorf("[reportDeviceLoginEvent] report login event failed, err:{%v}", err)
	}
}
