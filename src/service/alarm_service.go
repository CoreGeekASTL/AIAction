/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"GIDS/common/constants"
	"GIDS/common/logger"
	"Go-chassis-extend/api/ServiceComb/go-chassis/client/rest"
	"Go-chassis-extend/api/ServiceComb/go-chassis/core"

	manager "AlarmSDK_GO/api/alarmapi"
	base "AlarmSDK_GO/api/base"
)

const (
	ValuesLen              = 2
	TimePeriodInit         = 3
	TimePeriodClean        = 5
	RetryTimes             = 360
	AlarmId300010          = "300010"
	alarmSuppressThresholdMs = 10 * 60 * 1000
	alarmRetrySleepSeconds  = 10
)

// 定义告警事件结构体
type AlarmEvent struct {
	AlarmID      string                   // 告警ID
	EventMessage string                   // 事件信息
	Type         base.GenerateOrClearType // 事件类型
}

type AlarmService interface {
	SendAlarm(alarmID, EventMessage string)
	ClearAlarm(alarmID, EventMessage string)
}

type alarmServiceImpl struct {
	alarms       map[string]int64
	alarmManager base.CSPAlarmManager
}

var alarmService alarmServiceImpl
var alarmEventChanel chan AlarmEvent

const maxAlarmListLen = 999

// NewAlarmService 初始化告警服务
func NewAlarmService() AlarmService {
	return &alarmService
}

func init() {
	logger.Infof("init alarm function")
	alarmService = alarmServiceImpl{
		alarms:       make(map[string]int64),
		alarmManager: manager.CSPInitAlarmSDK(os.Getenv(constants.EnvAppId), constants.ServiceName, os.Getenv(constants.NODENAME), manager.GetNodeIP()),
	}
	alarmEventChanel = make(chan AlarmEvent, maxAlarmListLen)
	alarmService.alarmManager.RegisterRsetClear()
	go alarmService.handleEvent()
}

func (a *alarmServiceImpl) handleEvent() {
	for {
		select {
		case event, ok := <-alarmEventChanel:
			if !ok {
				break
			}
			if event.Type == base.GenerateAlarm {
				a.sendAlarm(event)
			} else if event.Type == base.ClearAlarm {
				a.clearAlarm(event)
			}
		}
	}
}

func (a *alarmServiceImpl) sendAlarm(event AlarmEvent) bool {
	lastSendTime, exists := a.alarms[event.AlarmID]
	// 如果存在且在10分钟内，则跳过
	now := time.Now().UnixMilli()
	if exists && now-lastSendTime < alarmSuppressThresholdMs {
		log.Println("An alarm was already reported within 10 minute; skipping this operation.")
		return true
	}
	success := a.reportAlarm(event)
	if success {
		a.alarms[event.AlarmID] = now
		logger.Infof("Alarm %s send successfully", event.AlarmID)
	} else {
		logger.Infof("Alarm %s send unsuccessfully", event.AlarmID)
	}
	return success
}

func (a *alarmServiceImpl) sendAlarmEvent(event AlarmEvent) {
	const threshold = 5
	select {
	case alarmEventChanel <- event:
	case <-time.After(time.Second * threshold):
		logger.Errorf("failed to report alarm evet:%v", event)
	}
}

// SendAlarm 发送告警
func (a *alarmServiceImpl) SendAlarm(alarmID, EventMessage string) {
	logger.Infof("enter SendAlarm function")
	alarmEvent := AlarmEvent{AlarmID: alarmID, EventMessage: EventMessage, Type: base.GenerateAlarm}
	a.sendAlarmEvent(alarmEvent)
}

// ClearAlarm 清除告警
func (a *alarmServiceImpl) ClearAlarm(alarmID, eventMessage string) {
	logger.Infof("enter ClearAlarm function")
	alarmEvent := AlarmEvent{AlarmID: alarmID, Type: base.ClearAlarm, EventMessage: eventMessage}
	a.sendAlarmEvent(alarmEvent)
}

func (a *alarmServiceImpl) clearAlarm(alarmEvent AlarmEvent) {
	if _, ok := a.alarms[alarmEvent.AlarmID]; !ok {
		logger.Infof("Alarm %s is not exist, skip it", alarmEvent.AlarmID)
	}
	success := a.reportAlarm(alarmEvent)
	if success {
		delete(a.alarms, alarmEvent.AlarmID)
		logger.Infof("Alarm %s clear successfully", alarmEvent.AlarmID)
	} else {
		logger.Errorf("Alarm %s clear unsuccessfully", alarmEvent.AlarmID)
	}
}

// genAlarm 生成告警对象
func (a *alarmServiceImpl) reportAlarm(alarmEvent AlarmEvent) bool {
	alarm := manager.InitCSPAlarm(alarmEvent.AlarmID, alarmEvent.Type)
	alarm.AppendParameter("kind", "App")
	alarm.AppendParameter("namespace", os.Getenv(constants.NAMESPACE))
	alarm.AppendParameter("sourceip", manager.GetNodeIP())
	alarm.AppendParameter("EventMessage", alarmEvent.EventMessage)
	alarm.AppendParameter("EventSource", os.Getenv(constants.SERVICENAME))
	alarm.AppendParameter("OriginalEventTime", time.Now().String())
	for j := 0; j < 2; j++ {
		result := a.alarmManager.SendAlarm(alarm)
		if result {
			logger.Infof("[AlarmServiceImpl] send alarm sccessfully")
			return true
		} else {
			time.Sleep(alarmRetrySleepSeconds * time.Second)
			logger.Errorf("[AlarmServiceImpl] send alarm fail, retry %d", j+1)
		}
	}
	// 重试失败后返回false
	logger.Errorf("[AlarmServiceImpl] send alarm failed after all retries")
	return false
}

// AlarmParamInfo 告警校验返回的 Parameters
type AlarmParamInfo struct {
	ParamName  string `json:"paramName"`  // 告警参数名
	ParamValue string `json:"paramValue"` // 告警参数值
}

// AlarmInfo 告警信息
type AlarmInfo struct {
	Location   string `json:"location,omitempty"`
	AppendInfo string `json:"appendInfo,omitempty"`
	AlarmId    string `json:"alarmId,omitempty"`
}

// AlarmResponse 告警响应
type AlarmResponse struct {
	Retdesc string      `json:"retdesc,omitempty"`
	Data    []AlarmInfo `json:"data,omitempty"`
	RetCode string      `json:"retcode,omitempty"`
}

// CleanAllActiveAlarm 清除活动告警
var AlarmList = []string{AlarmId300010}

func CleanAllActiveAlarm() bool {
	alarmIds := strings.Join(AlarmList, "&")
	alarmInfos, err := GetAllActiveAlarmFromFMService(alarmIds)
	// 解决升级场景FM服务未启动导致告警不清除问题，增加重试操作到半小时
	for i := 0; i < RetryTimes; i++ {
		if err == nil {
			break
		}
		time.Sleep(time.Second * TimePeriodClean)
		logger.Errorf("Get AlarmInfo err , and sleep 5s retry %d times", i)
		alarmInfos, err = GetAllActiveAlarmFromFMService(alarmIds)
	}
	if err != nil {
		logger.Errorf("Get AlarmInfo err ")
		return false
	}
	for alarmId, alarmInfo := range alarmInfos {
		for _, params := range alarmInfo {
			if (params.ParamName == "sourceip" || params.ParamName == "源ip") && params.ParamValue == manager.GetNodeIP() {
				clearHistoryAlarm(AlarmEvent{AlarmID: alarmId})
			}
		}
	}
	return true
}

func clearHistoryAlarm(alarmEvent AlarmEvent) {
	logger.Infof("enter ClearErrorAlarm: %+v", alarmEvent)
	alarmEvent.Type = base.ClearAlarm
	result := alarmService.reportAlarm(alarmEvent)
	if result {
		logger.Infof("alarm clear success ,alarmInfo is %v", alarmEvent)
	} else {
		logger.Errorf("alarm clear failed ,alarmInfo is %v", alarmEvent)
	}
}

var mutex sync.RWMutex

// GetAllActiveAlarmFromFMService 获取FM的活动告警
func GetAllActiveAlarmFromFMService(alarmIds string) (map[string][]AlarmParamInfo, error) {
	logger.Errorf("Enter GetAllActiveAlarmFromFMService")
	mutex.Lock()
	defer mutex.Unlock()

	// 函数返回类型定义
	var mapAlarmsRet = make(map[string][]AlarmParamInfo)
	var errRet error

	// 局部变量定义
	var mapAlarms = make(map[string][]AlarmParamInfo)
	cmd := GetActiveAlarms

	// 异步从FM获取所有的的告警
	ch := make(chan string)
	go func() {
		jsonParam := `{"cmd":"` + cmd + `","language":"en-us","data":{"appId":"` + os.Getenv("APPID") + `","alarmIds":"` +
			alarmIds + `"}}`
		requestParams := []byte(jsonParam)
		logger.Errorf("GetAllActiveAlarmFromFMService  parameter is : ", jsonParam)
		var micServiceName = "FMService"
		mapAlarms, errRet = handlerActivityAlarmData(micServiceName, requestParams)
		ch <- "success"
	}()
	select {
	case _, ok := <-ch:
		if !ok {
			logger.Errorf("channel closed, mapAlarmsRet is %v", mapAlarmsRet)
			return mapAlarmsRet, errRet
		}
		mapAlarmsRet = mapAlarms
		logger.Errorf("get alarms from fm ended in 3 seconds, mapAlarmsRet is %v, errRet is %v", mapAlarmsRet, errRet)
	case <-time.After(time.Second * TimePeriodInit):
		errRet = errors.New("get alarms from fm timeout")
		logger.Errorf("get alarms timeout, mapAlarmsRet is %v", mapAlarmsRet)
	}

	return mapAlarmsRet, errRet
}

// handlerActivityAlarmData 主动从FMService获取活动告警，并处理数据格式
func handlerActivityAlarmData(micServiceName string, requestParams []byte) (map[string][]AlarmParamInfo, error) {
	var mapAlarms = make(map[string][]AlarmParamInfo)
	var alarmResp AlarmResponse
	fmServiceURL := "/fmOperation/v1/alarms/get_alarms"
	bodyStrJSON, err := OSHttpsGetRequestByCSE(fmServiceURL, micServiceName, POST, requestParams)
	if err != nil {
		// 当获取body失败
		logger.Errorf("get body  from fmservice error")
		return mapAlarms, err
	}

	logger.Errorf("get body from fmservice is %s", bodyStrJSON)

	err = json.Unmarshal([]byte(bodyStrJSON), &alarmResp)
	if err != nil {
		logger.Errorf("unmarshal alarmInfo body from fmservice  failed!", err)
		return mapAlarms, err
	}
	logger.Errorf("get activity alarmData %s", alarmResp)
	if alarmResp.RetCode != ResultCode {
		logger.Errorf("can not get activity alarmData , reCode is %s",
			alarmResp.RetCode)
		return mapAlarms, errors.New("can not get activity alarmData")
	}
	for _, value := range alarmResp.Data {
		alarmId := value.AlarmId
		location := strings.Split(value.Location, ",")
		alarmParamArray := make([]AlarmParamInfo, 0)
		for _, value := range location {
			values := strings.Split(value, "=")
			if len(values) >= ValuesLen {
				alarmParam := getAlarmParam(alarmId, values[0], values[1])
				alarmParamArray = append(alarmParamArray, alarmParam)
			}
		}

		logger.Errorf("get alarmInterface from localCache fail alarmId is %v", alarmId)

		mapAlarms[alarmId] = alarmParamArray
	}
	return mapAlarms, nil
}

func getAlarmParam(alarmId, paramKey, paramValue string) AlarmParamInfo {
	logger.Errorf("get getAlarmParam from localCache fail alarmId is %v", alarmId)
	return AlarmParamInfo{ParamName: paramKey, ParamValue: paramValue}
}

// GetActiveAlarms 从FMServices获取告警
const (
	RespOK          = 200
	POST            = "POST"
	GetActiveAlarms = "GET_ACTIVE_ALARMS"
	ResultCode      = "0" // 0表示正常，非0表示有异常
)

// RetMsg 返回的结构信息
type RetMsg struct {
	Err        error
	StatusCode int
	Body       []byte
}

// OSHttpsGetRequestByCSE 功能：通过CSE对http的get方法进行封装，只需要传入需要的url参数，就可以返回body值
// 入参：希望获得的URL
// 返回值：1、执行成功还是失败；2、返回的body结果
func OSHttpsGetRequestByCSE(url string, microServiceName string, method string, body []byte) (string, error) {
	var bodyStr string
	logger.Debugf("request by cse method is %v  url is %v , "+
		"and params is %v ", method, "cse://"+microServiceName+"/"+strings.TrimLeft(url, "/"), string(body))
	request, err := rest.NewRequest(method, "cse://"+microServiceName+"/"+strings.TrimLeft(url, "/"), body)
	if err != nil {

		logger.Errorf("[OSHttpsGetRequestByCSE]NewRequest error !")
		return bodyStr, err
	}
	defer request.Close()
	logger.Debugf("request do.....")
	// 创建http链接
	response, err := core.NewRestInvoker().ContextDo(context.TODO(), request)
	if err != nil {
		logger.Errorf("[OSHttpsGetRequestByCSE]get vnfids from MK failed")
		return bodyStr, err
	}
	if response == nil {
		logger.Errorf("[OSHttpsGetRequestByCSE]response is null")
		return bodyStr, err
	}
	defer response.Close()

	if response.GetStatusCode() == RespOK {
		bodyStr = string(response.ReadBody())
		return bodyStr, nil
	}
	logger.Infof("Return code is %d", response.GetStatusCode())
	return bodyStr, errors.New("return code is not 200")
}
