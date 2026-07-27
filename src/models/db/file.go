// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package db

import (
	"fmt"
	"runtime"
	"strings"

	"GIDS/common/logger"

	"github.com/beego/beego/v2/client/orm"
)

type File struct {
	ID        int            `orm:"auto;pk;column(id)"`
	Bucket    string         `orm:"column(bucket)"`
	Name      string         `orm:"column(name)"`
	Content   ByteArrayField `orm:"column(content);type(bytea)"`
	Size      int64          `orm:"column(size)"`
	CreatedAt string         `orm:"column(created_at)"`
}

func (s *File) TableName() string {
	return "t_file"
}

type ByteArrayField []byte

// set value
func (e *ByteArrayField) SetRaw(value interface{}) error {
	if value == nil {
		return nil
	}
	switch d := value.(type) {
	case []byte:
		*e = d
	case string:
		*e = []byte(d)
	default:
		return fmt.Errorf("[ByteArrayField] unsupported type")
	}
	return nil
}

func (e *ByteArrayField) RawValue() interface{} {
	return *e
}

// FieldType specified type
func (e *ByteArrayField) FieldType() int {
	return orm.TypeTextField
}

func (f *ByteArrayField) String() string {
	return string(*f)
}

func getStackTrace() string {
	buf := make([]byte, 1024*10) // 10KB的缓冲区，足够存储一般堆栈信息
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func init() {
	// 使用defer和recover捕获异常
	defer func() {
		if err := recover(); err != nil {
			// 捕获到panic，进行处理
			logger.Errorf("UploadPluginPackage error: %v", err)
			stackTrace := getStackTrace()
			// 格式化输出，每行前缀添加"  "
			formattedTrace := strings.ReplaceAll(stackTrace, "\n", "\n  ")
			logger.Errorf("  %s\n", formattedTrace)
		}
	}()
	orm.RegisterModel(&File{})
}
