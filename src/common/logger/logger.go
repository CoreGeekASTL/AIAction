// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package logger
package logger

import (
	"fmt"

	"Go-chassis-extend/api/ServiceComb/go-chassis/core/lager"
)

// Infof log function handler
func Infof(format string, args ...interface{}) {
	lager.Logger.Infof(format, args...)
}

// Warnf log function handler
func Warnf(format string, args ...interface{}) {
	lager.Logger.Warnf(nil, format, args...)
}

// Debugf log function handler
func Debugf(format string, args ...interface{}) {
	lager.Logger.Debugf(format, args...)
}

// Errorf log function handler
func Errorf(format string, args ...interface{}) {
	lager.Logger.Errorf(nil, format, args...)
}

// TeeErrorf log function handler
func TeeErrorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args)
	lager.Logger.Errorf(nil, format, args...)
	return err
}

// Fatalf log function handler
func Fatalf(format string, args ...interface{}) {
	lager.Logger.Fatalf(nil, format, args...)
}
