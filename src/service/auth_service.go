// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package service

import (
	"regexp"
	"sync"

	"github.com/beego/beego/v2/client/orm"

	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
)

// whiteListDao 白名单存取接口，dao.WhiteListDao 隐式实现，单测注入 fake
type whiteListDao interface {
	Count() (int64, error)
	GetByIMEI(imei string) (*db.AuthWhitelist, error)
	InsertMulti(records []db.AuthWhitelist) error
	ClearAndInsert(records []db.AuthWhitelist) error
	ListAll() ([]db.AuthWhitelist, error)
}

// idPattern IMEI/IMSI 严格 15 位纯数字
var idPattern = regexp.MustCompile(`^[0-9]{15}$`)

// AuthService 终端联合鉴权服务
type AuthService interface {
	// Check 联合鉴权：返回是否放行、IMEI/IMSI 格式是否合法
	Check(imei, imsi string) (pass bool, formatValid bool)
	// ClearCache 清空鉴权缓存，白名单导入成功后调用
	ClearCache()
}

type authServiceImpl struct {
	store whiteListDao
	cache *authCache
}

var _ AuthService = &authServiceImpl{}

var (
	authServiceOnce     sync.Once
	authServiceInstance *authServiceImpl
)

func NewAuthService() AuthService {
	authServiceOnce.Do(func() {
		authServiceInstance = &authServiceImpl{
			store: dao.NewWhiteListDao(),
			cache: newAuthCache(),
		}
	})
	return authServiceInstance
}

// Check 联合鉴权：格式非法短路 → 缓存查询 → 逃生态判定 → 按 IMEI 查行比对 IMSI（组合精确匹配同一行）
func (s *authServiceImpl) Check(imei, imsi string) (bool, bool) {
	if !idPattern.MatchString(imei) || !idPattern.MatchString(imsi) {
		return false, false
	}
	key := imei + "_" + imsi
	if result, ok := s.cache.get(key); ok {
		return result, true
	}
	result := s.checkFromStore(imei, imsi)
	s.cache.set(key, result)
	return result, true
}

// ClearCache 清空缓存，白名单变更后立即生效
func (s *authServiceImpl) ClearCache() {
	s.cache.clear()
}

// checkFromStore 回源判定：表空逃生态放行；DB 异常 fail-open 放行并记错误日志，避免阻断主流程
func (s *authServiceImpl) checkFromStore(imei, imsi string) bool {
	count, err := s.store.Count()
	if err != nil {
		logger.Errorf("[auth] count white list failed, err: [%v], fail-open", err)
		return true
	}
	if count == 0 {
		return true
	}
	wl, err := s.store.GetByIMEI(imei)
	if err != nil {
		if err == orm.ErrNoRows {
			return false
		}
		logger.Errorf("[auth] get white list by imei failed, err: [%v], fail-open", err)
		return true
	}
	return wl.IMSI == imsi
}
