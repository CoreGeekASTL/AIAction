/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package main
package main

import (
	"fmt"
	"os"
	"time"

	report "CSPGSOMF/ModulekeeperSDK/modulekeeperapi"
	RunlogSDK "CSPGSOMF/RunlogSDK/runlogapi"
	"CSPGSOMF/TransportSDK/transportapi"
	ntp "CSPNTP_SDK_GO/api"
	gsfapi "Go-chassis-extend/api/GSF/api"
	gsfapibase "Go-chassis-extend/api/GSF/api/base"
	"Go-chassis-extend/api/GSF/sdk/log/logapi"
	gceBase "Go-chassis-extend/api/base"

	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/cert"
	"GIDS/common/conf"
	"GIDS/common/constants"
	"GIDS/common/cse"
	"GIDS/common/https"
	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/routers"
	"GIDS/scheduler"
	"GIDS/service"
	"GIDS/utils/flagutil"
)

// 常量定义
const (
	RETRYTIMES        = 60
	InitGSFRETRYTIMES = 360
	InitGSFSleepTime  = 5
	InitDMSSleepTime  = 30
	defaultHttpPort   = 40050
)

func main() {
	done := make(chan bool)

	logapi.AddFilterFileName("logger/logger.go")
	logger.Infof("******************BEGIN gids******************")

	c := conf.Instance()
	flagutil.Parse(c)
	conf.SetDefault(c)

	https.InitMuenClient()
	// 初始化gsf框架
	initGSF()
	https.Init()

	// 初始化TransportSDK
	transportapi.Init()

	ntp.Init()
	// 日志和mk初始化
	runLogInit()
	go dao.EnsureConnectGaussDB()
	// 注册服务实例
	registerInstance(done)

	cse.Init()
	service.StartRefreshConfigTask()

	// 注册内部服务
	if err := startInternalServer(); err != nil {
		logger.Errorf("failed to start internal server:%v,exit", err)
		os.Exit(1)
	}

	// 注册外部服务
	if err := startExternalHttpsServer(); err != nil {
		logger.Errorf("failed to start external server:%v,exit", err)
		os.Exit(1)
	}

	// 启动数据清理定时任务
	scheduler.StartDataCleanupScheduler()

	// 启动 csp 话统监控，用于上报运营指标
	monitorService := service.NewMonitorService()
	go monitorService.InitMonitorSchedule()

	go service.CleanAllActiveAlarm()
	const maxRetryTimes = 5
	cse.NewCse().Report(maxRetryTimes)
	logger.Infof("******************END SERVICE REGIS******************")
	<-done
}

// runLogInit 初始化runlog
func runLogInit() {
	RunlogSDK.InitServer()
	report.Init()
	report.AddServiceName("gids")
}

func gsfStartHandler(done chan<- bool, podName string) {
	if len(podName) == 0 {
		gsfapi.CspStart()
		return
	}
	gsfapi.CspStart(gsfapi.WithLocation(&gceBase.Location{PodName: podName, Essential: []string{"podName"}}))
	done <- true
}

// 注册服务实例
func registerInstance(done chan<- bool) {
	// CSP框架启动（会阻塞当前协程，需要的话用协程启动）
	logger.Infof("GSF Started ...")
	podName := os.Getenv("PODNAME")
	logger.Infof("get env podName: %s", podName)
	go gsfStartHandler(done, podName)
}

// 初始化注册
func initGSF() {
	var err error
	// ========CSP 框架初始化流程开始===============
	// 初始化CSP 框架 不成功 30s重试一次
	for i := 0; i < InitGSFRETRYTIMES; i++ {
		err = gsfapi.CspInit()
		if err != nil && i == (InitGSFRETRYTIMES-1) {
			logger.Errorf("gids gsfapi.CspInit retry time out, err: %v", err)
			time.Sleep(5 * time.Second) // todo test
			logger.Fatalf("%v", err)
		}
		if err != nil {
			time.Sleep(InitGSFSleepTime * time.Second)
			logger.Errorf("gids gsfapi.CspInit fail, err: %v", err)
		} else {
			logger.Infof("gids gsfapi.CspInit success.")
			break
		}
	}

	// 注册优雅退出回调
	handler := &GracefulExitHandler{}
	gsfapi.RegistExitHandler(handler)
	gsfapi.HealthCheckStart(gsfapibase.RestProtocal)
}

// 监听http服务
func startInternalServer() error {
	ip, err := https.GetLocalIP("FABRIC_ETH", "bond-base")
	if err != nil {
		return fmt.Errorf("get local ip failed err is %v", err)
	}
	port, err := beego.AppConfig.Int("httpport")
	if err != nil {
		logger.Errorf("failed to get http port:%v", err)
		port = 9090
	}

	internalServer := https.NewHttpServer(ip, port)
	routers.RegisterInternalRouter(internalServer)
	internalServer.Run()
	logger.Infof("http internal server running on port %v", port)
	return nil
}

// 监听外部服务地址
func startExternalHttpsServer() error {
	ip, err := https.GetLocalIP("SC_TRUNK_ETH", "bond-external")
	if err != nil {
		return fmt.Errorf("get local ip failed err is %v", err)
	}

	enableHttp := os.Getenv(constants.EnableHTTP)
	if enableHttp == "true" {
		// PORT环境变量优先级更高(40050)
		httpPort := beego.AppConfig.DefaultInt("moon::httpport", defaultHttpPort)
		logger.Infof("get env enableHttp value is %v , start http server", enableHttp)
		externalHttpServer := https.NewHttpServer(ip, httpPort)
		routers.RegisterExternalRouter(externalHttpServer)
		externalHttpServer.Run()
		logger.Infof("http external server running on port %v", httpPort)
		cse.NewCse().AddChainEndpoint(fmt.Sprintf("http://%s:%d", ip, httpPort))
	}

	// 设置的TLS_PORT环境变量优先级更高
	port, err := beego.AppConfig.Int("httpsport")
	if err != nil {
		logger.Errorf("failed to get https port:%v", err)
		port = 40051
	}

	externalServer := https.NewHttpsServer(ip, port)
	routers.RegisterExternalRouter(externalServer)
	externalServer.Run()
	logger.Infof("https external server running on port %v", port)

	cert.InitCert()
	err = cert.SubscribeCert(externalServer)
	if err != nil {
		logger.Fatalf("failed to SubscribeCert")
	}
	cse.NewCse().AddChainEndpoint(fmt.Sprintf("https://%s:%d", ip, port))
	return nil
}

// GracefulExitHandler 结构体
type GracefulExitHandler struct {
}

// Exit 优雅退出
func (gh *GracefulExitHandler) Exit() {
	logger.Infof("Executing exit handler, stopping scheduler...")
	scheduler.StopDataCleanupScheduler()
}
