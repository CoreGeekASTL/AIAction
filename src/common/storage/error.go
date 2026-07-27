// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."
package storage

import "errors"

const (
	ErrNotExist StorageError = "data not exist"
)

type StorageError string

func (s StorageError) Error() string {
	return string(s)
}

func IsNotExist(err error) bool {
	return errors.Is(err, ErrNotExist)
}
