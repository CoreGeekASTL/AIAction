/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package service
package service

import (
	"context"
	"errors"
	"time"

	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
	"github.com/beego/beego/v2/client/orm"
)

// ConfigCenterService config center interface
type ConfigCenterService interface {
	GetConfig(key string) (string, bool)
	InsertOrUpdateConfig(center db.ConfigCenter) error
	Refresh()
	GetFromDB(key string) (db.ConfigCenter, bool)
}

// RefreshInterval interval to refresh config center cache
const RefreshInterval = 5 * time.Minute

type configCenterServiceImpl struct {
	configs  map[string]string
	dao     *dao.ConfigCenterDao
	stopChan chan struct{}
}

// GetFromDB get config from DB
func (c *configCenterServiceImpl) GetFromDB(key string) (db.ConfigCenter, bool) {
	oldConfig := db.ConfigCenter{
		Key: key,
	}
	err := c.dao.Get(&oldConfig, "Key")
	if err != nil {
		return db.ConfigCenter{}, false
	}
	return oldConfig, true
}

// Refresh config cache refresh
func (c *configCenterServiceImpl) Refresh() {
	logger.Infof("start to refresh config center")
	var configList []db.ConfigCenter
	if err := c.dao.List(&configList); err != nil {
		logger.Errorf("[Refresh] list all config failed : %v", err)
		return
	}
	var configMap = make(map[string]string, 100)
	for _, config := range configList {
		configMap[config.Key] = config.Value
	}
	c.configs = configMap
}

// GetConfig get config from cache
func (c *configCenterServiceImpl) GetConfig(key string) (string, bool) {
	value, ok := configCenter.configs[key]
	return value, ok
}

// InsertOrUpdateConfig insert or update config
func (c *configCenterServiceImpl) InsertOrUpdateConfig(center db.ConfigCenter) error {
	err := c.dao.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		oldConfig := db.ConfigCenter{
			Key: center.Key,
		}
		err := c.dao.Get(&oldConfig, "Key")
		if errors.Is(err, orm.ErrNoRows) {
			center.UpdatedAt = time.Now().Format(time.DateTime)
			return c.dao.Insert(&center)
		} else if err != nil {
			return err
		}
		center.ID = oldConfig.ID
		center.UpdatedAt = time.Now().Format(time.DateTime)
		return c.dao.Update(&center)
	})
	return err
}

var configCenter *configCenterServiceImpl

func init() {
	configCenter = &configCenterServiceImpl{
		dao: dao.NewConfigCenterDao(),
	}
	configCenter.stopChan = make(chan struct{})
}

// NewConfigCenterService create ConfigCenterService
func NewConfigCenterService() ConfigCenterService {
	return configCenter
}

// StartRefreshConfigTask start timer task to refresh config center cache
func StartRefreshConfigTask() {
	logger.Infof("start to refresh config center interval")
	go func() {
		service := NewConfigCenterService()
		ticker := time.NewTicker(RefreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				service.Refresh()
			case <-configCenter.stopChan:
				logger.Infof("stopping config center refresh task")
				return
			}
		}
	}()
}

func StopRefreshConfigTask() {
	close(configCenter.stopChan)
}
