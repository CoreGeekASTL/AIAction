// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package controllers include web controllers
package controllers

import (
	"GIDS/common/constants/retcode"
	"GIDS/models/events"
	"GIDS/models/req"
	"GIDS/models/resp"
	"GIDS/service"
)

type EventController struct {
	BaseController
	eventService service.EventService
}

func (c *EventController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/app-api/center/public/client/sendClientEvent":      "POST:SendClientEvent",
			"/app-api/center/public/client/sendAppUseTimesEvent": "POST:SendAppUseTimesEvent",
		},
	}
}

func (c *EventController) Prepare() {
	c.eventService = service.NewEventService()
}

func (c *EventController) SendClientEvent() {
	request := new(req.ClientEventRequest)
	err := c.RequestBodyUnmarshalTo(request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}

	event := events.NewInfo(events.Client)
	event.SetEventData(events.ClientEventData{
		HSMan:   request.HSMan,
		HSType:  request.HSType,
		AppType: request.AppType,
		IMEI:    request.IMEI,
		IMSI:    request.IMSI,
		Type:    request.Type,
	})

	if err := c.eventService.ReportEvent(event); err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}

	response := resp.DataResponse{
		BaseResponse: resp.BaseResponse{
			Code:    retcode.Success,
			Message: "record success",
		},
		Data: true,
	}
	c.OK(response)
}

func (c *EventController) SendAppUseTimesEvent() {
	request := new(req.AppUseTimesEvent)
	err := c.RequestBodyUnmarshalTo(request)
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}

	event := events.NewInfo(events.AppUseTimes)
	event.SetEventData(events.AppUseTimesEvent{
		UseTimes: request.UseTimes,
		HSMan:    request.HSMan,
		HSType:   request.HSType,
		AppType:  request.AppType,
		AppId:    request.AppId,
		EXTType:  request.EXTType,
		PlayMode: request.PlayMode,
		SCWidth:  request.SCWidth,
		SCHeight: request.SCHeight,
		IMEI:     request.IMEI,
		IMSI:     request.IMSI,
	})

	if err := c.eventService.ReportEvent(event); err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}

	response := resp.DataResponse{
		BaseResponse: resp.BaseResponse{
			Code:    retcode.Success,
			Message: "record success",
		},
		Data: true,
	}
	c.OK(response)
}
