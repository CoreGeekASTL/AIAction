// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package https s for https server implementation
package https

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BeegoHttpServer struct {
	server *beego.HttpServer
	ip     string
	port   int
}

func (h *BeegoHttpServer) InsertFilter(pattern string, pos int, filter beego.FilterFunc, opts ...beego.FilterOpt) BeegoServer {
	h.server.InsertFilter(pattern, pos, filter, opts...)
	return h
}

func (h *BeegoHttpServer) Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) BeegoServer {
	h.server.Router(rootpath, c, mappingMethods...)
	return h
}

func (h *BeegoHttpServer) Run() {
	go h.server.Run("")
}

type BeegoServer interface {
	Run()
	Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) BeegoServer
	InsertFilter(pattern string, pos int, filter beego.FilterFunc, opts ...beego.FilterOpt) BeegoServer
}

func NewHttpServer(ip string, port int) *BeegoHttpServer {
	config := *beego.BConfig
	server := beego.NewHttpServerWithCfg(&config)
	server.Cfg.Listen.HTTPPort = port
	server.Cfg.Listen.HTTPAddr = ip
	return &BeegoHttpServer{
		server: server,
		ip:     ip,
		port:   port,
	}
}
