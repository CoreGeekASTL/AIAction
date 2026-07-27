// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package https s for https server implementation
package https

import (
	"fmt"
	"net"
	"os"

	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/logger"
)

const restartExitCode = 3

func NewHttpsServer(ip string, port int) *BeegoHttpsServer {
	return &BeegoHttpsServer{
		server:        newBeegoHttpsServer(ip, port),
		ip:            ip,
		port:          port,
		certInfo:      CertInfo{},
		restartChan:   make(chan CertInfo, 1),
		stopChan:      make(chan struct{}),
		isServerReady: false,
	}
}

// BeegoHttpsServer is a https server
type BeegoHttpsServer struct {
	server        *beego.HttpServer
	ip            string
	port          int
	certInfo      CertInfo
	restartChan   chan CertInfo
	stopChan      chan struct{}
	isServerReady bool
}

func (b *BeegoHttpsServer) close() {
	err := b.server.Server.Close()
	if err != nil {
		logger.Errorf("failed to close beego https server : %v", err)
	}
}

func newBeegoHttpsServer(ip string, port int) *beego.HttpServer {
	config := *beego.BConfig
	server := beego.NewHttpServerWithCfg(&config)
	server.Cfg.Listen.EnableHTTPS = true
	server.Cfg.Listen.HTTPSPort = port
	server.Cfg.Listen.EnableHTTP = false
	server.Cfg.Listen.HTTPSAddr = ip
	return server
}

func (b *BeegoHttpsServer) Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) BeegoServer {
	b.server.Router(rootpath, c, mappingMethods...)
	return b
}

func (b *BeegoHttpsServer) InsertFilter(pattern string, pos int, filter beego.FilterFunc, opts ...beego.FilterOpt) BeegoServer {
	b.server.InsertFilter(pattern, pos, filter, opts...)
	return b
}

// GetLocalIP 从环境变量获取网口名称，以及网口ip
// ethEnv 网口名称承载的环境变量名，当不存在时，使用网口名称
func GetLocalIP(ethEnv, defaultEth string) (string, error) {
	var localIP string
	var err error
	eth := os.Getenv(ethEnv)
	if len(eth) == 0 {
		eth = defaultEth
	}
	localIP, err = getEthIP(eth)
	if err != nil || len(localIP) == 0 {
		// 本地开发回退：找不到指定网口（如 Windows 上无 bond-base/bond-external）时，
		// 直接监听回环地址，保证 9090 等端口能正常启动。
		logger.Errorf("get %v local ip failed err is [%v], fallback to 127.0.0.1", eth, err)
		return "127.0.0.1", nil
	}
	return localIP, nil
}

func getEthIP(ethName string) (string, error) {
	iface, err := net.InterfaceByName(ethName)
	if err != nil {
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		switch ip := addr.(type) {
		case *net.IPNet:
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		default:
			logger.Warnf("Fail to get ipv4, types error")
		}
	}
	return "", fmt.Errorf("fail to get ipv4, no ip")
}

// Run start a beego server to listen
func (b *BeegoHttpsServer) Run() {
	go b.monitorCertificate()
}

func (b *BeegoHttpsServer) Stop() {
	close(b.stopChan)
}

func (b *BeegoHttpsServer) monitorCertificate() {
	for {
		logger.Infof("server ready to start, now wait for certificate upload")
		select {
		case certInfo := <-b.restartChan:
			b.updateCertInfo(certInfo)
			if !b.needStartServer() {
				logger.Infof("cert not upload, skip restart https server")
				continue
			}
			logger.Infof("certificate uploaded, try to start server")
			b.server.Server.TLSConfig = GetTLS(b.certInfo, ServerType)
			go b.server.Run("")
			logger.Infof("https server running on port %v", b.server.Cfg.Listen.HTTPSPort)
			b.isServerReady = true
		case <-b.stopChan:
			logger.Infof("stopping certificate monitor")
			return
		}
	}
}

func (b *BeegoHttpsServer) updateCertInfo(certInfo CertInfo) {
	if certInfo.KeyFile != "" {
		b.certInfo.CertFile = certInfo.CertFile
		b.certInfo.KeyFile = certInfo.KeyFile
		b.certInfo.KeyPwd = certInfo.KeyPwd
	}
	if certInfo.CaFile != "" {
		b.certInfo.CaFile = certInfo.CaFile
	}
}

func (b *BeegoHttpsServer) needStartServer() bool {
	// 证书已上传但服务端未启动，继续拉起端口
	if b.certInfo.KeyFile != "" && b.certInfo.CaFile != "" && !b.isServerReady {
		return true
	}
	// 证书已上传但服务端已启动，重启服务，更新证书
	if b.certInfo.KeyFile != "" && b.certInfo.CaFile != "" && b.isServerReady {
		logger.Infof("cert update, exit to restart gids try to update server cert")
		os.Exit(restartExitCode)
	}
	// 证书未获取成功，跳过此次事件
	return false
}

// UpdateCert start a beego server to listen
func (b *BeegoHttpsServer) UpdateCert(certInfo CertInfo) {
	logger.Infof("certificate uploaded, try to start server")
	b.restartChan <- certInfo
	return
}
