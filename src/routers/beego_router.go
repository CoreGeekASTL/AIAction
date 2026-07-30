/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package routers 定义Beego路由配置。
package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/https"
	"GIDS/common/logger"
	"GIDS/controllers"
)

// RegisterExternalRouter 注册外部路由
func RegisterExternalRouter(server https.BeegoServer) {
	// 限流过滤器
	server.InsertFilter("*", beego.BeforeRouter, controllers.OverLoadFilter)
	registerController(server, &controllers.ExLoginController{})
	registerController(server, &controllers.CacheController{})
	registerController(server, &controllers.EventController{})
	registerController(server, &controllers.ExFileController{})
	registerController(server, &controllers.TestController{})
}

// RegisterInternalRouter 注册内部路由
func RegisterInternalRouter(server https.BeegoServer) {
	// 限流过滤器
	server.InsertFilter("*", beego.BeforeRouter, controllers.OverLoadFilter)
	registerController(server, &controllers.CacheController{})
	registerController(server, &controllers.EventController{})
	registerController(server, &controllers.FileController{})
	registerController(server, &controllers.LoginController{})
	registerController(server, &controllers.ManagementController{})
	registerController(server, &controllers.PluginController{})
	registerController(server, &controllers.TrafficStatsController{})
	registerController(server, &controllers.ConfigCenterController{})
	registerController(server, &controllers.AuthController{})
}

func registerController(server https.BeegoServer, controller controllers.IController) {
	routeInfo := controller.RouteInfo()
	// 注册路由
	for k, v := range routeInfo.RouteMapping {
		logger.Infof("beego register router %v", k)
		server.Router(k, controller, v)
	}

	registerFilters(server, routeInfo, "")
}

func registerFilters(server https.BeegoServer, routeInfo controllers.RouteInfo, routePathPre string) {
	// 注册全局过滤器（匹配所有路由）
	for k, v := range routeInfo.Filters {
		var pos = beego.BeforeExec
		if k == controllers.After {
			pos = beego.AfterExec
		}
		server.InsertFilter("/*", pos, v, beego.WithReturnOnOutput(false))
	}

	// 注册带前缀的过滤器（匹配特定前缀下的所有子路由）
	if routePathPre != "" {
		for k, v := range routeInfo.Filters {
			var pos = beego.BeforeExec
			if k == controllers.After {
				pos = beego.AfterExec
			}
			server.InsertFilter(routePathPre+"/*", pos, v, beego.WithReturnOnOutput(false))
		}
	}
}
