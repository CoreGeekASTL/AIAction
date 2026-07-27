/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package controllers

import (
	"fmt"
	"net/http"

	"GIDS/common/logger"
	"code.huawei.com/fusionstage/greatwall-sdk-go/overloadcontroller"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

const FilterConfKey = "APIService"

var retryAfter = 3

// OverloadedResponse 过载提示
const OverloadedResponse = "Too many requests, please retry later."

func init() {
	err := overloadcontroller.Init()
	if err != nil {
		logger.Errorf("init greatwall failed: %v", err)
	}
}

func OverLoadFilter(ctx *beecontext.Context) {
	logger.Infof("OverLoadFilter start")
	dimNameValues := map[string]string{
		FilterConfKey: ctx.Request.URL.Path + "/" + ctx.Input.Method(),
	}
	logger.Infof("[OverLoadFilter] %s", dimNameValues[FilterConfKey])

	isGranted, err := overloadcontroller.Process(dimNameValues)
	if err != nil {
		logger.Errorf("overloadcontroller process failed: %v", err)
	}
	if !isGranted {
		ctx.ResponseWriter.Header().Add("Retry-After", fmt.Sprintf("%d", retryAfter))
		ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
		if _, err := ctx.ResponseWriter.Write([]byte(OverloadedResponse)); err != nil {
			logger.Errorf("OverLoadFilter write response error: %v", err)
		}
		return
	}
}
