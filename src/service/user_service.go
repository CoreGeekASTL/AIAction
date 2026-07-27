// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package service
package service

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"

	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
	"GIDS/models/req"
)

const (
	minDesensitizeLength = 4
	desensitizeKeepChars = 2
)

type UserService interface {
	CreateOrUpdateUser(request *req.LoginAuthRequest) error
	GetUserBind(sessionId string) (*db.UserBind, error)
	ExpiredUserBind(sessionId string) error
	UpdateUserBind(request *req.UpdateUserBindRequest) error
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{
		ud:  dao.NewUserDaoDao(),
		ubd: dao.NewUserBindDao(),
	}
}

var _ UserService = &UserServiceImpl{}

type UserServiceImpl struct {
	ud  *dao.UserDao
	ubd *dao.UserBindDao
}

func (u *UserServiceImpl) UpdateUserBind(request *req.UpdateUserBindRequest) error {
	ub, err := u.GetUserBind(request.SessionID)
	if err != nil {
		return err
	}

	if request.BrowserInstance != "" {
		ub.BrowserInstance = request.BrowserInstance
	}
	if request.InnerBrowserEndpoint != "" {
		ub.InnerBrowserEndpoint = request.InnerBrowserEndpoint
	}
	if request.InnerMediaEndpoint != "" {
		ub.InnerMediaEndpoint = request.InnerMediaEndpoint
	}
	if request.ControlEndpoint != "" {
		ub.ControlEndpoint = request.ControlEndpoint
	}
	if request.MediaEndpoint != "" {
		ub.MediaEndpoint = request.MediaEndpoint
	}
	if request.ControlTlsEndpoint != "" {
		ub.ControlTlsEndpoint = request.ControlTlsEndpoint
	}
	if request.MediaTlsEndpoint != "" {
		ub.MediaTlsEndpoint = request.MediaTlsEndpoint
	}

	ub.Heartbeats = time.Now().Format(time.DateTime)
	return u.ubd.Update(ub)
}

func (u *UserServiceImpl) ExpiredUserBind(sessionId string) error {
	ub, err := u.GetUserBind(sessionId)
	if err != nil {
		return err
	}

	ub.Heartbeats = time.Now().Format(time.DateTime)
	return u.ubd.Update(ub)
}

func (u *UserServiceImpl) GetUserBind(sessionId string) (*db.UserBind, error) {
	ub := &db.UserBind{
		Key: sessionId,
	}
	err := u.ubd.Get(ub)
	return ub, err
}

func (u *UserServiceImpl) CreateOrUpdateUser(request *req.LoginAuthRequest) error {
	newUser := getUserFormReq(request)
	user := &db.User{
		Key: fmt.Sprintf("%s_%s", request.IMEI, request.IMSI),
	}
	err := u.ud.Get(user)
	if err != nil && err != orm.ErrNoRows {
		logger.Errorf("[CreateOrUpdateUser] get User[%+v] failed, err: %v", desensitize(newUser.GetKey()), err)
		return fmt.Errorf("get user[%+v] failed: %v", desensitize(newUser.GetKey()), err)
	}
	if err == orm.ErrNoRows {
		user = newUser
		user.CreatedAt = time.Now().Format(time.DateTime)
		user.UpdatedAt = time.Now().Format(time.DateTime)
		err = u.ud.Insert(user)
	} else {
		user.UpdatedAt = time.Now().Format(time.DateTime)
		err = u.ud.Update(user)
	}
	if err != nil {
		logger.Errorf("[CreateOrUpdateUser] update User[%+v] failed, err: %v", desensitize(newUser.GetKey()), err)
		return err
	}
	return nil
}

func getUserFormReq(request *req.LoginAuthRequest) *db.User {
	var user = db.User{
		Key:          fmt.Sprintf("%s_%s", request.IMEI, request.IMSI),
		Manufacturer: request.Manufacturer,
		Model:        request.Model,
		ExtendModel:  request.ExtendModel,
		Width:        request.Width,
		Height:       request.Height,
		Country:      request.Country,
		Platform:     request.Platform,
		MCC:          request.MCC,
		MNC:          request.MNC,
		DeviceType:   request.DeviceType,
	}
	return &user
}

func desensitize(target string) string {
	if len(target) <= minDesensitizeLength {
		return target
	}

	return target[:desensitizeKeepChars] + "******" + target[len(target)-desensitizeKeepChars:]
}
