// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package service
package service

import (
	"sync"

	"GIDS/common/conf"
	"GIDS/common/event"
	"GIDS/common/logger"
	"GIDS/models/events"
)

type EventService interface {
	ReportEvent(event *events.Info) error
}

func NewEventService() *EventServiceImpl {
	once.Do(initEventStorageFactory)

	var err error
	eventStorage, err := event.GetFactory().Get(DefaultEventStorage)
	if err != nil {
		logger.Errorf("illegal eventStorage [%s], will use Default EventStorage", DefaultEventStorage)
		eventStorage, _ = event.GetFactory().Get(DefaultEventStorage)
	}
	return &EventServiceImpl{
		eventStorage: eventStorage,
	}
}

type EventServiceImpl struct {
	eventStorage event.Storage
}

var once = sync.Once{}

const (
	DefaultEventStorage = "localAuditComponent"
)

func initEventStorageFactory() {
	event.InitFactory()
	event.GetFactory().Register(DefaultEventStorage, event.NewLocalEventStorage(DefaultEventStorage, conf.Instance().Logger.EventFile))
}

func (e *EventServiceImpl) ReportEvent(event *events.Info) error {
	return e.eventStorage.Record(event)
}
