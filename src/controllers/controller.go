// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package controllers include web controllers
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"

	"GIDS/common/constants/retcode"
	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/req"
	"GIDS/models/resp"
)

const (
	Before FilterAction = "before"
	After  FilterAction = "after"
)

type FilterAction string

type RouteInfo struct {
	RouteMapping map[string]string
	Filters      map[FilterAction]beego.FilterFunc
}

type IController interface {
	beego.ControllerInterface
	RouteInfo() RouteInfo
}

type BaseController struct {
	beego.Controller
}

func (c *BaseController) RouteInfo() RouteInfo {
	return RouteInfo{}
}

func (c *BaseController) QueryParameter(name string) string {
	return c.Ctx.Input.Query(name)
}

func (c *BaseController) PathParameter(name string) string {
	return c.Ctx.Input.Param(name)
}

func (c *BaseController) Body() io.ReadCloser {
	return c.Ctx.Request.Body
}

func (c *BaseController) AddHeader(header, value string) {
	c.Ctx.Output.Header(header, value)
}

func (c *BaseController) ResponseWriter() http.ResponseWriter {
	return c.Ctx.ResponseWriter
}

func (c *BaseController) Request() *http.Request {
	return c.Ctx.Request
}

func (c *BaseController) RequestBodyUnmarshalTo(param req.IRequest) error {
	inputBody, err := io.ReadAll(c.Body())
	if err != nil {
		return logger.TeeErrorf("read rquest body failed, err: %v", err)
	}
	logger.Infof("request body is %s", string(inputBody))
	err = json.Unmarshal(inputBody, param)
	if err != nil {
		errInfo := fmt.Sprintf("json unmarshal req failed, err: %v", err)
		logger.Errorf(errInfo)
		return errors.New(errInfo)
	}
	err = param.Validate()
	if err != nil {
		errInfo := fmt.Sprintf("req param valid failed, err: %v", err)
		logger.Errorf(errInfo)
		return errors.New(errInfo)
	}
	return nil
}

func (c *BaseController) OK(data interface{}) {
	if data == nil {
		data = resp.BaseResponse{
			Code:    retcode.Success,
			Message: "success",
		}
	}
	err := c.writeHeaderAndJSON(http.StatusOK, data, "application/json")
	if err != nil {
		logger.Errorf("return output %v failed, turn to InternalServiceError", data)
		c.InternalServiceError()
	}
	respData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("marshal response data failed: %v", err)
	} else {
		logger.Infof("response is %s", string(respData))
	}
}

func (c *BaseController) Failed(data resp.BaseResponse) {
	err := c.writeHeaderAndJSON(http.StatusBadRequest, data, "application/json")
	if err != nil {
		logger.Errorf("return output %v failed, turn to InternalServiceError", data)
		c.InternalServiceError()
	}
}

func (c *BaseController) NotFound() {
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
}

func (c *BaseController) InternalServiceError() {
	data := resp.BaseResponse{
		Code:    retcode.InternalFailed,
		Message: "InternalServiceError",
	}
	err := c.writeHeaderAndJSON(http.StatusInternalServerError, data, "application/json")
	if err != nil {
		logger.Errorf("return output %v failed", data)
	}
}

// write marshalls the value to JSON and set the Content-Type Header.
func (c *BaseController) writeHeaderAndJSON(status int, v interface{}, contentType string) error {
	if v == nil {
		c.Ctx.ResponseWriter.WriteHeader(status)
		// do not write a nil representation
		return nil
	}
	c.AddHeader("Content-Type", contentType)
	c.Ctx.ResponseWriter.WriteHeader(status)
	return json.NewEncoder(c.Ctx.ResponseWriter).Encode(v)
}

func (c *BaseController) insertData(tag string, param req.ITable, d dao.BaseInterface) {
	err := c.RequestBodyUnmarshalTo(param)
	if err != nil {
		logger.Errorf("[%v] unmarshal failed, err: [%v], request: [%v]", tag, err, param)
		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}

	err = d.Insert(param)
	if err != nil {
		logger.Errorf("[%v] create failed, err is [%v]", tag, err)
		c.Failed(resp.BaseResponse{Code: retcode.InternalFailed, Message: err.Error()})
		return
	}
	c.OK(nil)
}
