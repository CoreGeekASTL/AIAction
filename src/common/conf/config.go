// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package config handle config
package conf

import (
	beego "github.com/beego/beego/v2/server/web"
	"os"
	"strings"
)

var (
	defaultRedisEndpoint        = beego.AppConfig.DefaultString("redis::endpoint", "")
	defaultNodeExternalEndpoint = beego.AppConfig.DefaultString("node::endpoint", "")
	defaultNodeHttpsEndpoint    = beego.AppConfig.DefaultString("node::httpsendpoint", "")
	defaultOSSEndpoint          = beego.AppConfig.DefaultString("oss::endpoint", "")
)

var config *Config

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	config = &Config{
		Logger: LoggerConfig{
			LogLevel: "INFO",
		},
		Redis: RedisConfig{Endpoint: defaultRedisEndpoint},
		Node: NodeConfig{
			HostName:         hostname,
			ExternalEndpoint: defaultNodeExternalEndpoint,
			HttpsEndpoint:    defaultNodeHttpsEndpoint,
		},
		OSS: OSSConfig{
			Endpoint:  defaultOSSEndpoint,
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			Token:     "",
		},
	}
}

type NodeConfig struct {
	ExternalEndpoint string `flag:"endpoint"`
	HttpsEndpoint    string `flag:"httpsendpoint"`
	HostName         string `flag:"hostname"`
}

// Config define start config
type Config struct {
	Logger LoggerConfig `flag:"log"`
	Redis  RedisConfig  `flag:"redis"`
	OSS    OSSConfig    `flag:"oss"`
	Node   NodeConfig   `flag:"node"`
}

// LoggerConfig define logger config
type LoggerConfig struct {
	LogFile   string `flag:"file" desc:"log file path"`
	LogLevel  string `flag:"level" desc:"log level: INFO/WARN/DEBUG"`
	EventFile string `flag:"event" desc:"file path of event log"`
}

type RedisConfig struct {
	Endpoint string `flag:"endpoint"`
	DB       int    `flag:"db" desc:"redis db number: 0-15"`
}

type OSSConfig struct {
	Endpoint  string `flag:"endpoint"`
	AccessKey string `flag:"accessKey"`
	SecretKey string `flag:"secretKey"`
	Token     string `flag:"token"`
}

// Instance get config instance
func Instance() *Config {
	return config
}

// SetDefault 针对尼日局点设置默认值
func SetDefault(c *Config) {
	if strings.Contains(c.Node.HttpsEndpoint, ">>") {
		c.Node.HttpsEndpoint = defaultNodeHttpsEndpoint
	}
}
