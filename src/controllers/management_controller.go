// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"GIDS/service"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/https"
	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
	"GIDS/models/resp"
)

const hoursPerDay = 24
const defaultRetryCount = 2

var moonConfig = "moon"
var interval = hoursPerDay * time.Hour

type BrowserConfig struct {
	RouteAPPConfigList []db.RouterAPPConfig `json:"routeAppConfigList,omitempty"`
	ChromeConfigList   []db.ChromeConfig    `json:"chromeConfigList,omitempty"`
	URLConfigs         []db.URLConfig       `json:"urlConfigList,omitempty"`
}

type ManagementController struct {
	BaseController
	cd           *dao.ConfigDao
	alarm        service.AlarmService
	configCenter service.ConfigCenterService
}

func (c *ManagementController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			// TODO 分布式多实例时，同步操作需要更改, 避免多实例去同步配置, 先获取分布式锁, 待增加定时任务
			"/rpc-api/center/config/syncBrowserConfig": "POST:SyncBrowserConfig",
			"/config/v1": "GET:ListConfig",
		},
	}
}

func (c *ManagementController) Prepare() {
	c.cd = dao.NewConfigDao()
	c.alarm = service.NewAlarmService()
	c.configCenter = service.NewConfigCenterService()
}

// 定期同步配置
func (c *ManagementController) updateConfigIfNeed() *db.Config {
	cfg := &db.Config{
		Type: moonConfig,
	}
	err := c.cd.Get(cfg, "Type")
	if err != nil {
		logger.Errorf("get config failed: %v", err)
		if err := c.syncBrowserConfig(); err != nil {
			logger.Errorf("sysnc browser config failed: %v", err)
		}
		return nil
	}
	updatedTime, err := time.Parse(time.DateTime, cfg.UpdatedAt)
	if err != nil {
		logger.Errorf("parse time failed: %v", err)
		if err := c.syncBrowserConfig(); err != nil {
			logger.Errorf("sysnc browser config failed: %v", err)
		}
		return nil
	}
	if time.Since(updatedTime) > interval {
		if err := c.syncBrowserConfig(); err != nil {
			logger.Errorf("sysnc browser config failed: %v", err)
		}
	}
	return cfg
}

func (c *ManagementController) ListConfig() {
	cfg := c.updateConfigIfNeed()
	if cfg == nil {
		cfg = &db.Config{
			Type: moonConfig,
		}
		if err := c.cd.Get(cfg, "Type"); err != nil {
			if errors.Is(err, orm.ErrNoRows) {
				c.NotFound()
			}
			logger.Errorf("get config failed: %v", err)
			c.InternalServiceError()
			return
		}
	}
	bc := &BrowserConfig{}
	if err := json.Unmarshal([]byte(cfg.Content), bc); err != nil {
		logger.Errorf("parse config failed: %v", err)
		c.InternalServiceError()
		return
	}
	c.OK(bc)
}

func (c *ManagementController) getMoonConfigUrl() (bool, string) {
	cfgUrl := beego.AppConfig.DefaultString("moon::configEndpoint", "")
	config, ok := c.configCenter.GetConfig("moon::configEndpoint")
	if ok && config != "" {
		cfgUrl = config
	}

	cfgHttpsUrl := beego.AppConfig.DefaultString("moon::httpsConfigEndpoint", "")
	httpsConfig, ok := c.configCenter.GetConfig("moon::httpsConfigEndpoint")
	if ok && httpsConfig != "" {
		cfgHttpsUrl = httpsConfig
	}

	enableHttps := beego.AppConfig.DefaultString("moon::enableHttps", "")
	enableHttpsConfig, ok := c.configCenter.GetConfig("moon::enableHttps")
	if ok && enableHttpsConfig != "" {
		enableHttps = enableHttpsConfig
	}

	if enableHttps == "true" && cfgHttpsUrl != "" {
		return true, cfgHttpsUrl
	}

	return false, cfgUrl
}

func (c *ManagementController) syncBrowserConfig() error {
	var (
		enableHttps, cfgUrl = c.getMoonConfigUrl()
		err                 error
		bc                  = &BrowserConfig{}
	)

	client := https.Instance()
	if enableHttps {
		client = https.MuenInstance()
	}
	logger.Infof("use cloud service url %v, enableHttps %v", cfgUrl, enableHttps)
	moonResp := https.NewRequest(client).
		Method(http.MethodGet).WithRetry(defaultRetryCount).
		URL(cfgUrl).Complete().Do()

	logger.Infof("[syncBrowserConfig] response: %v", moonResp.Response())
	if !moonResp.IsSuccessCode() || moonResp.Error() != nil {
		logger.Errorf("sync browser config failed, status is %d, err is %v", moonResp.StatusCode(), moonResp.Error())
		return moonResp.Error()
	}

	if err = moonResp.ResponseToStruct(&resp.DataResponse{Data: bc}); err != nil {
		logger.Errorf("parse sync browser config response failed:%v", err)
		return err
	}

	contentB, err := json.Marshal(bc)
	if err != nil {
		logger.Errorf("parse sync browser config response failed:%v", err)
		return err
	}
	content := string(contentB)
	cfg := &db.Config{
		Type:    moonConfig,
		Content: content,
	}
	if err := c.insertOrUpdate(cfg); err != nil {
		logger.Errorf("sync browser config failed, err is %v", err)
		return err
	}
	return nil
}

func (c *ManagementController) SyncBrowserConfig() {
	logger.Infof("start SyncBrowserConfig")
	if err := c.syncBrowserConfig(); err != nil {
		c.InternalServiceError()
		// 上报告警
		c.alarm.SendAlarm(service.AlarmId300010, "Failed to sync browser configuration")
		return
	}
	// 恢复告警
	c.alarm.ClearAlarm(service.AlarmId300010, "Browser configuration synced successfully")
	logger.Infof("End SyncBrowserConfig")
	c.OK(nil)
}

func (c *ManagementController) insertOrUpdate(cfg *db.Config) error {
	oldCfg := &db.Config{
		Type: cfg.Type,
	}
	err := c.cd.Get(oldCfg, "Type")
	if err != nil && err != orm.ErrNoRows {
		return err
	}
	if err == orm.ErrNoRows {
		cfg.CreatedAt = time.Now().Format(time.DateTime)
		cfg.UpdatedAt = time.Now().Format(time.DateTime)
		return c.cd.Insert(cfg)
	}
	oldCfg.Content = cfg.Content
	oldCfg.UpdatedAt = time.Now().Format(time.DateTime)
	return c.cd.Update(oldCfg)
}
