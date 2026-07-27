/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package cert 提供证书信息处理和TLS配置相关功能。
package cert

import (
	"CSPGSOMF/CertSDK/api/base"
	"CSPGSOMF/CertSDK/api/certapi"

	"GIDS/common/https"
	"GIDS/common/logger"
)

var exCertMgr base.CSPExCertManager
var externalInfos []base.CspExSceneInfo
var externalCaInfos []base.CspExSceneInfo
var serverInfos []base.CspExSceneInfo

var newCertInfo = https.CertInfo{}

func InitCert() {
	logger.Infof("start subscribe sbg certificate scene")
	err := certapi.CspCertSDKInit()
	if err != nil {
		logger.Fatalf("init cert sdk failed: %v", err)
	}
	exCertMgr = certapi.GetExCertManagerInstance()
}

func InitCertScene() error {
	sceneCaInfo := base.CspExSceneInfo{
		SceneName:   "sbg_external_ca_certificate",
		SceneDescCN: "云浏览器外部CA证书",
		SceneDescEN: "SBG external CA certificate",
		SceneType:   1,
		Feature:     0,
	}
	sceneDeviceInfo := base.CspExSceneInfo{
		SceneName:   "sbg_external_device_certificate",
		SceneDescCN: "云浏览器外部设备证书",
		SceneDescEN: "SBG external Device Certificate",
		SceneType:   2,
		Feature:     0,
	}

	serverCaInfo := base.CspExSceneInfo{
		SceneName:   "sbg_server_ca_certificate",
		SceneDescCN: "云浏览器服务端CA证书",
		SceneDescEN: "SBG server CA certificate",
		SceneType:   1,
		Feature:     0,
	}
	serverDeviceInfo := base.CspExSceneInfo{
		SceneName:   "sbg_server_device_certificate",
		SceneDescCN: "云浏览器服务端设备证书",
		SceneDescEN: "SBG server Device Certificate",
		SceneType:   2,
		Feature:     0,
	}
	externalCaInfos = append(externalCaInfos, sceneCaInfo)
	externalInfos = append(externalInfos, sceneDeviceInfo)
	serverInfos = append(serverInfos, serverCaInfo)
	serverInfos = append(serverInfos, serverDeviceInfo)

	// 换包重启时SubscribeExCert不会清理可能存在的旧订阅, 所以需要单独处理取消旧订阅场景，避免客户端更新时服务端联动复位
	err := exCertMgr.UnsubscribeExCert("gids-muen", []base.CspExSceneInfo{sceneCaInfo, serverCaInfo, serverDeviceInfo})
	if err != nil {
		return err
	}
	err = exCertMgr.UnsubscribeExCert("gids-muenCa", []base.CspExSceneInfo{sceneDeviceInfo, serverCaInfo, serverDeviceInfo})
	if err != nil {
		return err
	}
	err = exCertMgr.UnsubscribeExCert("gids", []base.CspExSceneInfo{sceneCaInfo, sceneDeviceInfo})
	if err != nil {
		return err
	}
	return nil
}

// SubscribeCert 订阅证书更新
func SubscribeCert(server *https.BeegoHttpsServer) error {
	err := InitCertScene()
	if err != nil {
		return err
	}
	err = exCertMgr.SubscribeExCert("gids-muen", externalInfos, exCertInfoHandler, "/opt/csp/gids/")
	if err != nil {
		return err
	}
	err = exCertMgr.SubscribeExCert("gids-muenCa", externalCaInfos, exCertInfoHandler, "/opt/csp/gids/")
	if err != nil {
		return err
	}
	err = exCertMgr.SubscribeExCert("gids", serverInfos, func(certInfo []*base.CspExCertInfo, notifyType int) error {
		return serverCertInfoHandler(server, certInfo, notifyType)
	}, "/opt/csp/gids/")
	if err != nil {
		return err
	}
	return nil
}

func exCertInfoHandler(certInfo []*base.CspExCertInfo, notifyType int) error {
	logger.Infof("get sbg external cert update, try update client")
	for _, info := range certInfo {
		res, err := exCertMgr.GetExCertPathInfo(info.SceneName)
		if err != nil {
			logger.Errorf("get cert path failed: %v", err)
		}
		pwd, err := exCertMgr.GetExCertPrivateKeyPwd(info.SceneName)
		if err != nil {
			logger.Errorf("get cert pwd failed: %v", err)
		}
		switch info.SceneName {
		case "sbg_external_ca_certificate":
			newCertInfo.CaFile = res.ExCaFilePath
		case "sbg_external_device_certificate":
			newCertInfo.CertFile = res.ExDeviceFilePath
			newCertInfo.KeyFile = res.ExPrivateKeyFilePath
			newCertInfo.KeyPwd = pwd
		default:
			logger.Infof("unknown external cert scene: %s", info.SceneName)
		}
	}
	https.MuenCertUpdate(newCertInfo)
	return nil
}

func serverCertInfoHandler(server *https.BeegoHttpsServer, certInfo []*base.CspExCertInfo, notifyType int) error {
	logger.Infof("get sbg server cert update, try update server")
	cert := https.CertInfo{}
	for _, info := range certInfo {
		res, err := exCertMgr.GetExCertPathInfo(info.SceneName)
		if err != nil {
			logger.Errorf("get cert path failed: %v", err)
		}
		pwd, err := exCertMgr.GetExCertPrivateKeyPwd(info.SceneName)
		if err != nil {
			logger.Errorf("get cert pwd failed: %v", err)
		}
		switch info.SceneName {
		case "sbg_server_ca_certificate":
			cert.CaFile = res.ExCaFilePath
		case "sbg_server_device_certificate":
			cert.CertFile = res.ExDeviceFilePath
			cert.KeyFile = res.ExPrivateKeyFilePath
			cert.KeyPwd = pwd
		default:
			logger.Infof("unknown server cert scene: %s", info.SceneName)
		}
	}
	server.UpdateCert(cert)
	return nil
}
