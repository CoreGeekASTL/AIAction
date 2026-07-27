/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package event 提供事件存储与查询功能。
package event

import (
	"fmt"
	"sync"

	"GIDS/common/logger"
	"GIDS/models/events"
)

type Storage interface {
	Record(logInfo *events.Info) error
}

// StorageFactory registry for location-storages
type StorageFactory interface {
	// Get storage by given location
	Get(location string) (Storage, error)
	// Register storage by given location
	Register(location string, storage Storage)
}

type eventStorageFactory struct {
	*sync.RWMutex
	factory map[string]Storage
}

var (
	factoryInstance StorageFactory
)

// InitFactory initialize empty factory
func InitFactory() {
	factoryInstance = &eventStorageFactory{
		RWMutex: &sync.RWMutex{},
		factory: map[string]Storage{},
	}
}

// GetFactory gets StorageFactory
func GetFactory() StorageFactory {
	return factoryInstance
}

func (sf *eventStorageFactory) Get(location string) (Storage, error) {
	sf.RLock()
	st, ok := sf.factory[location]
	sf.RUnlock()
	if !ok {
		return nil, fmt.Errorf("eventStorage %s not found in factory", location)
	}
	if st == nil {
		return nil, fmt.Errorf("nil eventStorage %s found in factory", location)
	}
	return st, nil
}

func (sf *eventStorageFactory) Register(location string, storage Storage) {
	if len(location) == 0 || storage == nil {
		logger.Errorf("failed to register eventStorageFactory because location or storage is nil")
		return
	}

	sf.Lock()
	sf.factory[location] = storage
	sf.Unlock()
}
