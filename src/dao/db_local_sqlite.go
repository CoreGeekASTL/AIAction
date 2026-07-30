/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.
 */

package dao

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/logger"

	_ "modernc.org/sqlite"
)

const localSqliteInitSql = `
CREATE TABLE IF NOT EXISTS t_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    "type" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    created_at TEXT DEFAULT '',
    updated_at TEXT DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_file (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bucket TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "content" BLOB,
    "size" INTEGER,
    created_at TEXT DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_plugin_package (
    "key" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "version" TEXT NOT NULL,
    package_name TEXT NOT NULL,
    plugin_type TEXT NOT NULL,
    bucket TEXT NOT NULL,
    active_status TEXT NOT NULL,
    if_active INTEGER NOT NULL DEFAULT 0,
    progress INTEGER NOT NULL,
    created_at TEXT DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_user (
    "key" TEXT NOT NULL PRIMARY KEY,
    manufacturer TEXT NOT NULL,
    model TEXT NOT NULL,
    extend_model TEXT DEFAULT '',
    country TEXT DEFAULT '',
    platform TEXT DEFAULT '',
    width TEXT DEFAULT '',
    height TEXT DEFAULT '',
    mcc TEXT DEFAULT '',
    mnc TEXT DEFAULT '',
    device_type TEXT DEFAULT '',
    created_at TEXT DEFAULT '',
    updated_at TEXT DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_user_bind (
    "key" TEXT NOT NULL PRIMARY KEY,
    browser_instance TEXT DEFAULT '',
    media_endpoint TEXT DEFAULT '',
    control_endpoint TEXT DEFAULT '',
    media_tls_endpoint TEXT DEFAULT '',
    control_tls_endpoint TEXT DEFAULT '',
    inner_media_endpoint TEXT DEFAULT '',
    inner_browser_endpoint TEXT DEFAULT '',
    "token" TEXT DEFAULT '',
    updated_at TEXT DEFAULT ''
);
CREATE TABLE IF NOT EXISTS t_media_traffic_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    app_type INTEGER NOT NULL,
    started_at TEXT DEFAULT '',
    finished_at TEXT DEFAULT '',
    out_bytes INTEGER NOT NULL,
    access_type INTEGER
);
CREATE INDEX IF NOT EXISTS idx_mtraffic_session ON "t_media_traffic_stats" (session_id);
CREATE INDEX IF NOT EXISTS idx_mtraffic_app ON "t_media_traffic_stats" (app_type);
CREATE INDEX IF NOT EXISTS idx_mtraffic_session_app ON "t_media_traffic_stats" (session_id,app_type);
CREATE TABLE IF NOT EXISTS t_control_traffic_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    app_type INTEGER NOT NULL,
    started_at TEXT DEFAULT '',
    finished_at TEXT DEFAULT '',
    out_bytes INTEGER NOT NULL,
    access_type INTEGER
);
CREATE INDEX IF NOT EXISTS idx_ctraffic_session ON "t_control_traffic_stats" (session_id);
CREATE INDEX IF NOT EXISTS idx_ctraffic_app ON "t_control_traffic_stats" (app_type);
CREATE INDEX IF NOT EXISTS idx_ctraffic_session_app ON "t_control_traffic_stats" (session_id,app_type);
CREATE TABLE IF NOT EXISTS t_session_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    app_type INTEGER NOT NULL,
    started_at TEXT DEFAULT '',
    finished_at TEXT DEFAULT '',
    tcp_unique_id TEXT
);
CREATE INDEX IF NOT EXISTS idx_session_session ON "t_session_stats" (session_id);
CREATE INDEX IF NOT EXISTS idx_session_app ON "t_session_stats" (app_type);
CREATE INDEX IF NOT EXISTS idx_session_session_app ON "t_session_stats" (session_id,app_type);
CREATE UNIQUE INDEX IF NOT EXISTS t_tcp_unique_id ON t_session_stats(tcp_unique_id);
CREATE TABLE IF NOT EXISTS t_config_center (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_key TEXT NOT NULL,
    config_value TEXT NOT NULL,
    config_describe TEXT DEFAULT '',
    enable INTEGER DEFAULT 1,
    updated_at TEXT DEFAULT ''
);
CREATE UNIQUE INDEX IF NOT EXISTS t_configs_key ON t_config_center(config_key);
CREATE TABLE IF NOT EXISTS t_white_list (
    imei TEXT NOT NULL PRIMARY KEY,
    imsi TEXT NOT NULL,
    created_at TEXT DEFAULT ''
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_white_list_imei_imsi ON t_white_list(imei, imsi);
`

var sqliteDriverOnce sync.Once

func initLocalSQLite() error {
	dbPath := beego.AppConfig.DefaultString("local::sqlitepath", "./data/gids.db")

	if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create sqlite data dir %s failed: %w", dir, err)
		}
	}

	var driverErr error
	sqliteDriverOnce.Do(func() {
		if probe, err := sql.Open("sqlite3", ":memory:"); err == nil {
			_ = probe.Close()
			return
		}
		base, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			driverErr = fmt.Errorf("open sqlite driver failed: %w", err)
			return
		}
		defer base.Close()
		sql.Register("sqlite3", base.Driver())
	})
	if driverErr != nil {
		return driverErr
	}

	if err := orm.RegisterDataBase("default", "sqlite3", dbPath); err != nil {
		if !strings.Contains(err.Error(), "already registered") {
			return fmt.Errorf("register sqlite database failed: %w", err)
		}
		logger.Infof("[initLocalSQLite] database alias default already registered, reuse it")
	}

	ormer = orm.NewOrmUsingDB("default")

	var firstErr error
	statements := strings.Split(localSqliteInitSql, ";")
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		logger.Infof("[initLocalSQLite] exec sql: %s", stmt)
		if _, err := ormer.Raw(stmt).Exec(); err != nil {
			logger.Errorf("[initLocalSQLite] exec sql#%d failed: %v", i, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	if firstErr != nil {
		return fmt.Errorf("init sqlite tables failed, first error: %w", firstErr)
	}
	logger.Infof("[initLocalSQLite] init local sqlite success, db=%s", dbPath)
	return nil
}
