// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package service
package service

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"GIDS/common/https"
	"GIDS/models/browsergateway"

	"GIDS/common/cse"
	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"

	"GIDS/common/conf"
	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
	"GIDS/models/req"
	"GIDS/models/resp"
)

const (
	defaultTTLMultiplier        = 2
	defaultHeartbeatsMultiplier = 3
)

var (
	defaultTTL        = defaultTTLMultiplier * time.Minute
	defaultHeartbeats = defaultHeartbeatsMultiplier * time.Minute
)

type BrowserService interface {
	RouteToInstance(request *req.LoginAuthRequest) (resp.LoginInfo, error)
	GetAllServiceInstances() []browsergateway.ServiceInstance
	GetAllReadyServiceInstances() []browsergateway.ServiceInstance
	UpdateUserToken(loginInfo *resp.LoginInfo, request *req.LoginAuthRequest) error
	PreOpenBrowser(request *req.LoginAuthRequest)
}

var _ BrowserService = &BrowserServiceImpl{}

func NewBrowserService() BrowserService {
	return &BrowserServiceImpl{
		ubd:        dao.NewUserBindDao(),
		cse:        cse.NewCse(),
		httpClient: https.Instance(),
	}
}

type BrowserServiceImpl struct {
	ubd        *dao.UserBindDao
	cse        cse.Cse
	httpClient https.HTTPDoer
}

func (b *BrowserServiceImpl) PreOpenBrowser(request *req.LoginAuthRequest) {
	initBrowserRequest := browsergateway.InitBrowserRequest{
		Factory:        request.Manufacturer,
		DevType:        request.Model,
		ExType:         request.ExtendModel,
		PlatType:       request.Platform,
		LcdWidth:       request.Width,
		LcdHeight:      request.Height,
		IMSI:           request.IMSI,
		IMEI:           request.IMEI,
		DeviceType:     request.DeviceType,
		ClientLanguage: request.ClientLanguage,
	}

	instances := b.GetAllReadyServiceInstances()
	for _, instance := range instances {
		go instancePreOpenBrowser(instance, initBrowserRequest, b.httpClient)
	}

}

func instancePreOpenBrowser(
	instance browsergateway.ServiceInstance,
	request browsergateway.InitBrowserRequest,
	httpClient https.HTTPDoer,
) {
	response := https.NewRequest(httpClient).
		Method("POST").
		URL(fmt.Sprintf("http://%s/browsergw/browser/preOpen", instance.BrowserInnerEndpoint)).
		ParamFromInterface(request).
		Complete().
		Do()

	if response.Response() == nil {
		logger.Errorf("BrowserGW %s pre open browser failed,err is %v",
			instance.BrowserInnerEndpoint, response.Error())
		return
	}
	if !response.IsSuccessCode() || response.Error() != nil {
		logger.Errorf("BrowserGW %s pre open browser, status is %d, err is %v",
			instance.BrowserInnerEndpoint, response.StatusCode(), response.Error())
		return
	}
	logger.Infof("BrowserGW %s pre open browser success, user %s",
		instance.BrowserInnerEndpoint, request.IMEI+"_"+request.IMSI)
}

func (b *BrowserServiceImpl) insertOrUpdate(new *db.UserBind) error {
	old := &db.UserBind{
		Key: new.Key,
	}
	err := b.ubd.Get(old, "Key")
	if err != nil && err != orm.ErrNoRows {
		return err
	}
	if err == orm.ErrNoRows {
		return b.ubd.Insert(new)
	}
	old.BrowserInstance = new.BrowserInstance
	old.ControlEndpoint = new.ControlEndpoint
	old.MediaEndpoint = new.MediaEndpoint
	old.ControlTlsEndpoint = new.ControlTlsEndpoint
	old.MediaTlsEndpoint = new.MediaTlsEndpoint
	old.InnerMediaEndpoint = new.InnerMediaEndpoint
	old.InnerBrowserEndpoint = new.InnerBrowserEndpoint
	old.Token = new.Token
	old.Heartbeats = new.Heartbeats
	return b.ubd.Update(old)
}

func (b *BrowserServiceImpl) reRouteToInstance(userBind *db.UserBind) (resp.LoginInfo, error) {
	err := b.routeToInstance(userBind)
	if err != nil {
		return resp.LoginInfo{}, err
	}
	createToken(userBind)
	if err = b.insertOrUpdate(userBind); err != nil {
		logger.Errorf("[BrowserServiceImpl] failed to save userBind [%s], err: [%v]", userBind.GetField(), err)
		return resp.LoginInfo{}, err
	}
	logger.Infof("[BrowserServiceImpl] success to create UserBind[%s],", userBind.GetField())
	return tranUserBindToLoginInfo(userBind), nil
}

func (b *BrowserServiceImpl) RouteToInstance(request *req.LoginAuthRequest) (resp.LoginInfo, error) {
	var userBind = &db.UserBind{
		Key: fmt.Sprintf("%s_%s", request.IMEI, request.IMSI),
	}

	// 1. 先判断UserBind是否存在或者有效，
	err := b.ubd.Get(userBind)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return resp.LoginInfo{}, err
	}
	if errors.Is(err, orm.ErrNoRows) || !b.validateUserBind(userBind) {
		// 2. 如果未分配userBind 或者userBind过期 或者实例异常
		return b.reRouteToInstance(userBind)
	}

	// 3. 如果成功获取UserBind, 直接返回
	logger.Infof("[BrowserServiceImpl] success to get UserBind[%s] from redis,", userBind.GetField())
	return tranUserBindToLoginInfo(userBind), nil
}

func (b *BrowserServiceImpl) checkIfUBExpired(userBind *db.UserBind) bool {
	heartBeats, err := time.Parse(time.DateTime, userBind.Heartbeats)
	if err != nil {
		logger.Errorf("[checkIfUBExpired] failed to parse heart beats: %v", err)
		return true
	}
	return heartBeats.Add(defaultHeartbeats).Before(time.Now())
}

func (b *BrowserServiceImpl) validateUserBind(userBind *db.UserBind) bool {
	bindInstance, ok := b.cse.GetBrowserGateWayInstanceByInnerEndpoint(userBind.BrowserInstance)
	if !ok {
		return false
	}
	return bindInstance.IsHealthy && !b.checkIfUBExpired(userBind)
}

func tranUserBindToLoginInfo(userBind *db.UserBind) resp.LoginInfo {
	return resp.LoginInfo{
		AuthInfo: resp.AuthInfo{
			Token:       userBind.Token,
			ExpiresTime: time.Now().Unix(),
			TimeAxis:    time.Now().Unix(),
		},
		AssignInfo: resp.AssignInfo{
			TcpAddr:             userBind.ControlEndpoint,
			TlsTcpAddr:          userBind.ControlTlsEndpoint,
			VideoMode:           1,
			ShortAddr:           conf.Instance().Node.ExternalEndpoint,
			NodeGateWayURL:      conf.Instance().Node.ExternalEndpoint,
			HttpsNodeGateWayUrl: conf.Instance().Node.HttpsEndpoint,
			HttpsShortAddr:      conf.Instance().Node.HttpsEndpoint,
			NodeIntranetWayURL:  userBind.InnerBrowserEndpoint,
			NodeCapacity:        userBind.BrowserCap,
		},
	}
}

func createToken(u *db.UserBind) {
	uid := uuid.New()
	u.Token = uid.String()
}

func (b *BrowserServiceImpl) routeToInstance(userBind *db.UserBind) error {
	browserInstance, err := b.assignInstance()
	if err != nil {
		return err
	}

	userBind.BrowserInstance = browserInstance.GetKey()
	userBind.MediaEndpoint = browserInstance.MediaExtendEndpoint
	userBind.InnerMediaEndpoint = browserInstance.MediaInnerEndpoint
	userBind.ControlEndpoint = browserInstance.ControlExtendEndpoint
	userBind.MediaTlsEndpoint = browserInstance.MediaTlsExtendEndpoint
	userBind.ControlTlsEndpoint = browserInstance.ControlTlsExtendEndpoint
	userBind.InnerBrowserEndpoint = browserInstance.BrowserInnerEndpoint
	userBind.BrowserCap = browserInstance.Cap
	return nil
}

func (b *BrowserServiceImpl) assignInstance() (browsergateway.ServiceInstance, error) {
	instances := b.GetAllReadyServiceInstances()
	sort.Sort(browsergateway.ServiceInstanceList(instances))
	// 优先分配给空闲的实例
	if len(instances) == 0 || instances[0].Used >= instances[0].Cap {
		return browsergateway.ServiceInstance{}, errors.New("no idle instances available")
	}
	return instances[0], nil
}

func (b *BrowserServiceImpl) GetAllReadyServiceInstances() []browsergateway.ServiceInstance {
	allSil := b.GetAllServiceInstances()
	var sil []browsergateway.ServiceInstance
	for i := range allSil {
		if allSil[i].PluginStatus == db.Complete && allSil[i].Cap > 0 && allSil[i].IsHealthy {
			sil = append(sil, allSil[i])
		}
	}
	return sil
}

func (b *BrowserServiceImpl) GetAllServiceInstances() []browsergateway.ServiceInstance {
	sil := b.cse.GetAllBrowserGateWayInstances()
	logger.Infof("[BrowserServiceImpl] finish getAllServiceInstances, instances:{%v}", len(sil))
	return sil
}

func (b *BrowserServiceImpl) UpdateUserToken(loginInfo *resp.LoginInfo, request *req.LoginAuthRequest) error {
	var ub = &db.UserBind{
		Key: fmt.Sprintf("%s_%s", request.IMEI, request.IMSI),
	}

	err := b.ubd.Get(ub)
	if err != nil {
		return err
	}

	ub.Token = loginInfo.Token
	ub.Heartbeats = time.Now().Format(time.DateTime)
	return b.ubd.Update(ub)
}
