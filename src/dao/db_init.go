//go:build !test

/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.
 */

package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/cse"
	"GIDS/common/https"
	"GIDS/common/logger"
	_ "GIDS/db/driver"
)

const initSql = `
CREATE TABLE IF NOT EXISTS t_config (
    id SERIAL PRIMARY KEY,  -- 自增主键，用于唯一标识每条记录
    "type" varchar(255) NOT NULL,
	"content" text NOT NULL,
	created_at varchar(255) DEFAULT '',
	updated_at varchar(255) DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_file (
    id SERIAL PRIMARY KEY,  -- 自增主键，用于唯一标识每条记录
	bucket varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	"content" bytea,
	"size" int8,
	created_at varchar(255) DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_plugin_package (
	"key" varchar(255) NOT NULL PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	"version" varchar(255) NOT NULL,
	package_name varchar(255) NOT NULL,
	plugin_type varchar(255) NOT NULL,
	bucket varchar(255) NOT NULL,
	active_status varchar(255) NOT NULL,
    if_active boolean NOT NULL DEFAULT false,
	progress int4 NOT NULL,
	created_at varchar(255) DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_user (
    "key" VARCHAR(128) NOT NULL PRIMARY KEY,  -- 以imsi作为主键，非空
    manufacturer VARCHAR(128) NOT NULL,  -- 制造商，默认空字符串
    model VARCHAR(128) NOT NULL,         -- 设备型号，默认空字符串
    extend_model VARCHAR(128) DEFAULT '',           -- 扩展型号，默认空字符串
    country VARCHAR(64) DEFAULT '',                 -- 国家/地区，默认空字符串
    platform VARCHAR(64) DEFAULT '',                -- 平台信息，默认空字符串
    width VARCHAR(32) DEFAULT '',                   -- 宽度，默认空字符串
    height VARCHAR(32) DEFAULT '',                  -- 高度，默认空字符串
    mcc VARCHAR(8) DEFAULT '',                      -- 移动国家代码，默认空字符串
    mnc VARCHAR(8) DEFAULT '',                      -- 移动网络代码，默认空字符串
    device_type VARCHAR(32) DEFAULT '',             -- 设备类型，默认空字符串
    created_at varchar(255) DEFAULT '',
	updated_at varchar(255) DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_user_bind (
    "key" VARCHAR(128) NOT NULL PRIMARY KEY,  -- 主键，对应结构体的Key字段
    browser_instance VARCHAR(255) DEFAULT '',  -- 浏览器实例标识，默认空字符串
    media_endpoint VARCHAR(255) DEFAULT '',    -- 媒体端点，默认空字符串
    control_endpoint VARCHAR(255) DEFAULT '',  -- 控制端点，默认空字符串
    media_tls_endpoint VARCHAR(255) DEFAULT '',    -- 媒体端点，默认空字符串
    control_tls_endpoint VARCHAR(255) DEFAULT '',  -- 控制端点，默认空字符串
    inner_media_endpoint VARCHAR(255) DEFAULT '',  -- 内部媒体端点，默认空字符串
    inner_browser_endpoint VARCHAR(255) DEFAULT '',  -- 内部浏览器端点，默认空字符串
    "token" VARCHAR(255) DEFAULT '',  -- 令牌，默认空字符串
    updated_at varchar(255) DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_media_traffic_stats (
    id SERIAL PRIMARY KEY,  -- 自增主键，用于唯一标识每条记录
	session_id varchar(255) NOT NULL, -- 用户信息
    app_type INT NOT NULL,  -- 应用类型
    started_at varchar(255) DEFAULT '', -- 连接开始时间
    finished_at varchar(255) DEFAULT '', -- 连接结束时间
	out_bytes BIGINT NOT NULL  -- 流量
);
CREATE INDEX IF NOT EXISTS idx_mtraffic_session ON "t_media_traffic_stats" (session_id);
CREATE INDEX IF NOT EXISTS idx_mtraffic_app ON "t_media_traffic_stats" (app_type);
CREATE INDEX IF NOT EXISTS idx_mtraffic_session_app ON "t_media_traffic_stats" (session_id,app_type);

CREATE TABLE IF NOT EXISTS t_control_traffic_stats (
    id SERIAL PRIMARY KEY,  -- 自增主键，用于唯一标识每条记录
	session_id varchar(255) NOT NULL, -- 用户信息
    app_type INT NOT NULL,  -- 应用类型
    started_at varchar(255) DEFAULT '', -- 连接开始时间
    finished_at varchar(255) DEFAULT '', -- 连接结束时间
	out_bytes BIGINT NOT NULL  -- 流量
);
CREATE INDEX IF NOT EXISTS idx_ctraffic_session ON "t_control_traffic_stats" (session_id);
CREATE INDEX IF NOT EXISTS idx_ctraffic_app ON "t_control_traffic_stats" (app_type);
CREATE INDEX IF NOT EXISTS idx_ctraffic_session_app ON "t_control_traffic_stats" (session_id,app_type);

CREATE TABLE IF NOT EXISTS t_session_stats (
    id SERIAL PRIMARY KEY,  -- 自增主键，用于唯一标识每条记录
	session_id varchar(255) NOT NULL, -- 用户信息
    app_type INT NOT NULL,  -- 应用类型
    started_at varchar(255) DEFAULT '' -- 连接开始时间
);
CREATE INDEX IF NOT EXISTS idx_session_session ON "t_session_stats" (session_id);
CREATE INDEX IF NOT EXISTS idx_session_app ON "t_session_stats" (app_type);
CREATE INDEX IF NOT EXISTS idx_session_session_app ON "t_session_stats" (session_id,app_type);
CREATE TABLE IF NOT EXISTS t_config_center (
    id SERIAL PRIMARY KEY,  -- 自增主键，用于唯一标识每条记录
	config_key varchar(255) NOT NULL, -- 配置名称
    config_value varchar(255) NOT NULL,  -- 配置值
    config_describe varchar(255) DEFAULT '', -- 配置描述
    enable boolean Default true, -- 配置是否启用
    updated_at varchar(255) DEFAULT '' -- 更新时间
);
CREATE UNIQUE INDEX IF NOT EXISTS t_configs_key ON t_config_center(config_key);
CREATE TABLE IF NOT EXISTS t_white_list (
    imei char(15) NOT NULL PRIMARY KEY,  -- 设备标识，15位纯数字
    imsi char(15) NOT NULL,  -- 用户身份标识，15位纯数字，Beego不支持复合pk，唯一性由联合索引兜底
    created_at varchar(255) DEFAULT ''
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_white_list_imei_imsi ON t_white_list(imei, imsi);
ALTER TABLE t_user_bind ADD COLUMN IF NOT EXISTS control_tls_endpoint VARCHAR(255);
ALTER TABLE t_user_bind ADD COLUMN IF NOT EXISTS media_tls_endpoint VARCHAR(255);

-- 添加finished_at和tcp_unique_id列
ALTER TABLE t_session_stats ADD COLUMN IF NOT EXISTS finished_at VARCHAR(255) DEFAULT '';
ALTER TABLE t_session_stats ADD COLUMN IF NOT EXISTS tcp_unique_id VARCHAR(36);
CREATE UNIQUE INDEX IF NOT EXISTS t_tcp_unique_id ON t_session_stats(tcp_unique_id);

-- 添加接入类型access_type
ALTER TABLE t_media_traffic_stats ADD COLUMN IF NOT EXISTS access_type INT;
ALTER TABLE t_control_traffic_stats ADD COLUMN IF NOT EXISTS access_type INT;
`
const interval = 5 * time.Second
const maxHealthCheckCount = 3

type dataSource struct {
	host  string
	alais string
}

type dbConnection struct {
	sync.Mutex
	currentDataSource dataSource            // 当前使用的数据源
	dataSourceMap     map[string]dataSource // 已建立连接的数据源
	countHealCheck    int                   // 数据源健康检查失败计数
}

var dbConnections *dbConnection
var dbServiceName = "" // db 服务名
var dbName = ""        // db名称

func (db *dbConnection) healthCheck() bool {
	db.Lock()
	defer db.Unlock()

	currentDB, err := orm.GetDB(db.currentDataSource.alais)
	if err != nil {
		logger.Errorf("get db failed: %v", err)
		return false
	}
	err = currentDB.Ping()
	if err != nil {
		logger.Errorf("faile to ping db: %v", err)
		db.countHealCheck++
	} else {
		db.countHealCheck = 0
	}
	return db.countHealCheck < maxHealthCheckCount
}

func (db *dbConnection) getCurrentDataSource() dataSource {
	db.Lock()
	defer db.Unlock()
	return db.currentDataSource
}

func (db *dbConnection) switchToAnotherDB(host, port, aliasName string) error {
	db.Lock()
	defer db.Unlock()
	if target, ok := db.dataSourceMap[host]; ok {
		db.currentDataSource = target
		ormer = orm.NewOrmUsingDB(target.alais)
		logger.Errorf("there exit connections to %s:%s, no need to reconnect", host, port)
		return nil
	}

	dataSourceUrl := getDataSourceUrl(host, port)
	logger.Errorf("s00893267 switchToAnotherDB dataSourceUrl is %s", dataSourceUrl)
	err := orm.RegisterDataBase(aliasName, "gaussdb_1", dataSourceUrl)
	dataSourceUrl = ""
	if err != nil {
		return err
	}
	source := dataSource{host: host, alais: aliasName}
	db.dataSourceMap[host] = source
	db.currentDataSource = source
	ormer = orm.NewOrmUsingDB(aliasName)
	logger.Infof("[ConnectGaussDB] connect database on %s:%s success", host, port)
	return nil
}

func init() {
	dbConnections = &dbConnection{
		currentDataSource: dataSource{"", ""},
		dataSourceMap:     make(map[string]dataSource),
		countHealCheck:    0,
	}

	dbServiceName = os.Getenv("DB_SERVICE_NAME")
	if dbServiceName == "" {
		dbServiceName = beego.AppConfig.DefaultString("gaussdb::servicename", "")
	}
	dbName = os.Getenv("DB_NAME")
	// 环境变量可能取到的是空或者bond-fabric
	if dbName == "" || dbName == "bond-fabric" {
		dbName = beego.AppConfig.DefaultString("gaussdb::gaussdbdbname", "")
		logger.Errorf("s00893267 [init] gaussdb::gaussdbdbname is %s", dbName)
	}
	logger.Errorf("s00893267 [init] dbName is %s", dbName)
}

func EnsureConnectGaussDB() {
	orm.BootStrap()
	// 注册数据库驱动
	err := orm.RegisterDriver("gaussdb_1", orm.DRPostgres)
	if err != nil {
		logger.Errorf("register driver err: %v", err)
		return
	}
	// 本地训战模式：免外部数据库，使用嵌入式 SQLite
	if os.Getenv("LOCAL_MODE") == "true" {
		if e := initLocalSQLite(); e != nil {
			logger.Errorf("init local sqlite failed: %v", e)
		}
		return
	}
	for {
		logger.Infof("connect gaussdb...")
		ip, port := getGaussDBIP()
		if ip == "" {
			time.Sleep(interval)
			continue
		}

		if err := dbConnections.switchToAnotherDB(ip, port, "default"); err != nil {
			logger.Errorf("failed to switch to new DB %s: %v", ip, err)
			time.Sleep(interval)
			continue
		}
		if err := initTables(); err != nil {
			logger.Errorf("failed to execute init sql: %v", err)
			continue
		}
		break
	}
	go checkDBStatus()
}

func checkDBStatus() {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			newDBIp, _ := getGaussDBIP()
			currentDataSource := dbConnections.getCurrentDataSource()
			if newDBIp == currentDataSource.host && dbConnections.healthCheck() {
				continue
			}
			logger.Warnf("[checkDBStatus] there's somethine with DB connection of %s, the new ip of DB is %s",
				currentDataSource.host, newDBIp)
			refresh()
		}
	}
}

func refresh() {
	logger.Infof("[refresh] start to refresh the DB connection")
	for {
		ip, post := getGaussDBIP()
		if ip == "" {
			logger.Warnf("[refresh] there is no master DB")
			time.Sleep(interval)
			continue
		}
		currentDataSource := dbConnections.getCurrentDataSource()
		if ip == currentDataSource.host {
			logger.Infof("[refresh] DB(%s) didn't change, no need to reconnect", ip)
			return
		}
		logger.Infof("[refresh] reconnect db from %s to %s", currentDataSource.host, ip)
		if err := dbConnections.switchToAnotherDB(ip, post, ip); err != nil {
			logger.Errorf("[refresh] failed to connect to %s : %v", ip, err)
			time.Sleep(interval)
			continue
		}
		return
	}
}

func initTables() error {
	initSqls := strings.Split(initSql, ";")
	for i := range initSqls {
		logger.Infof("exec sql: %s", initSqls[i])
		_, err := ormer.Raw(initSqls[i]).Exec()
		if err != nil {
			logger.Errorf("init tables err: %v", err)
			return err
		}
	}
	logger.Infof("init tables success")
	return nil
}

type gaussDbInfo struct {
	ServiceName  string `json:"SERVICE_NAME"`
	MasterIpAddr string `json:"MASTER_IP_ADDR"`
	Port         string `json:"PORT"`
	NewPwd       string `json:"NEW_DB_PWD"`
	DbName       string `json:"DB_NAME"`
	DbUser       string `json:"DB_USER"`
}

func getDataSourceUrl(host, port string) string {
	dataSourceUrl, err := getDataSourceFromDBService(host, port, dbServiceName, dbName)
	if err != nil || dataSourceUrl == "" {
		logger.Errorf("failed to get db info from service: %v", err)
		if dbServiceName == "GaussDB" {
			dataSourceUrl = getDataSourceByConfig(host)
			logger.Errorf("s00893267 getDataSourceUrl dataSourceUrl after getconf")
		}
	}
	return dataSourceUrl
}

// 读取默认配置
func getDataSourceByConfig(host string) string {
	user := beego.AppConfig.DefaultString("gaussdb::gaussdbuser", "")
	dbname := beego.AppConfig.DefaultString("gaussdb::gaussdbdbname", "")
	port := beego.AppConfig.DefaultString("gaussdb::gaussdbport", "")
	password := beego.AppConfig.DefaultString("gaussdb::databasepassword", "")
	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname
	return connStr
}

func getDataSourceFromDBService(host, port, serviceName, dbName string) (string, error) {
	dbUrl := fmt.Sprintf("https://%s:%s/service/api/getGaussdbInfor?serviceName=%s&dbName=%s", host, port,
		serviceName, dbName)
	logger.Errorf("s00893267 db service host is %s, dbName is %s, port is %s, serviceName is %s", host, dbName, port, serviceName)
	instance := https.InnerInstance()
	response := https.NewRequest(instance).URL(dbUrl).Method("GET").Complete().Do()
	if !response.IsSuccessCode() || response.Error() != nil {
		logger.Errorf("failed to get db info, url is:%s, status is %d, err is %v", dbUrl,
			response.StatusCode(), response.Error())
		return "", errors.New("failed to get gaussDbInfo")
	}
	body, err := io.ReadAll(response.ResponseBody())
	// 解析db service返回的数据
	logger.Errorf("s00893267 db service err is %s", err)

	bodyStr := strings.Replace(string(body), `#`, `"`, -1)
	bodyStr = strings.Trim(bodyStr, `"`)
	info := &gaussDbInfo{}
	logger.Errorf("s00893267 gaussDbInfo ServiceName is %s, MasterIpAddr is %s, Port is %s, NewPwd is %s, DbName is %s, DbUser is %s ", info.ServiceName, info.MasterIpAddr, info.Port, info.NewPwd, info.DbName, info.DbUser)
	err = json.Unmarshal([]byte(bodyStr), info)
	if err != nil {
		return "", err
	}
	connStr := "host=" + info.MasterIpAddr + " port=" + info.Port + " user=" + info.DbUser + " password=" + info.
		NewPwd + " dbname=" + info.DbName
	logger.Errorf("s00893267 connStr is %s", connStr)

	info = nil
	return connStr, nil
}

func getGaussDBIP() (string, string) {
	instances, err := cse.NewCse().GetAllMicroServiceInstanceInfo(dbServiceName)
	if err != nil {
		logger.Errorf("failed to get ndpdb intance form cse, err: %v", err)
		return "", ""
	}
	endpoints := make(map[string]struct{})
	for _, instance := range instances {
		if instance.Status != "UP" {
			continue
		}
		if instance.Properties["status"] != "M" { // 只连主gauss
			continue
		}
		for _, endpoint := range instance.EndpointsList {
			endpoints[endpoint] = struct{}{}
		}
	}

	for ep := range endpoints {
		ip, port, err := extractIP(ep)
		if err != nil {
			continue
		}
		return ip, port
	}
	return "", ""
}

func extractIP(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		logger.Errorf("failed to parse endpoint %s, err: %v", endpoint, err)
		return "", "", err
	}
	if u.Host == "" {
		return "", "", logger.TeeErrorf("failed to get host in endpoint %s", endpoint)
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		logger.Errorf("failed to split host and point in endpoint %s, err: %v", endpoint, err)
	}
	return host, port, nil
}
