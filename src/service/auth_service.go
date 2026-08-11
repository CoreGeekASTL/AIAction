// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package service
package service

import (
	"errors"
	"regexp"
	"sync"

	"github.com/beego/beego/v2/client/orm"

	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
)

// imeiPattern IMEI/IMSI 必须是 15 位纯数字
var imeiPattern = regexp.MustCompile(`^[0-9]{15}$`)

// whiteListDaoInterface 白名单 DAO 接口，便于测试替换
type whiteListDaoInterface interface {
	Count() (int64, error)
	GetByIMEIAndIMSI(imei, imsi string) error
	InsertMulti(list *[]db.WhiteList) error
	ClearAndInsert(list *[]db.WhiteList) error
	ListAll(list *[]db.WhiteList) error
}

// AuthService 终端联合鉴权服务
type AuthService interface {
	AuthIMEI(imei, imsi string) (allowed bool, formatValid bool)
}

type authServiceImpl struct {
	whiteListDao whiteListDaoInterface
	cache        *authCache
}

var (
	authServiceInstance AuthService
	authServiceOnce     sync.Once
)

// NewAuthService 获取鉴权服务单例
func NewAuthService() AuthService {
	authServiceOnce.Do(func() {
		authServiceInstance = &authServiceImpl{
			whiteListDao: dao.NewWhiteListDao(),
			cache:        getAuthCache(),
		}
	})
	return authServiceInstance
}

// authCacheKey 构造 IMEI+IMSI 联合缓存键
func authCacheKey(imei, imsi string) string {
	return imei + "|" + imsi
}

// AuthIMEI 联合鉴权：格式校验短路 → 缓存查询 → 逃生态判定 → 联合精确匹配，DB 异常安全优先拒绝
func (s *authServiceImpl) AuthIMEI(imei, imsi string) (bool, bool) {
	if !imeiPattern.MatchString(imei) || !imeiPattern.MatchString(imsi) {
		return false, false
	}
	key := authCacheKey(imei, imsi)
	if result, hit := s.cache.get(key); hit {
		return result, true
	}
	authImportLock.RLock()
	defer authImportLock.RUnlock()
	count, err := s.whiteListDao.Count()
	if err != nil {
		logger.Errorf("[AuthIMEI] count white list failed, err: [%v]", err)
		return false, true
	}
	if count == 0 {
		s.cache.set(key, true)
		return true, true
	}
	err = s.whiteListDao.GetByIMEIAndIMSI(imei, imsi)
	if err == nil {
		s.cache.set(key, true)
		return true, true
	}
	if errors.Is(err, orm.ErrNoRows) {
		s.cache.set(key, false)
		return false, true
	}
	logger.Errorf("[AuthIMEI] query white list failed, err: [%v]", err)
	return false, true
}
