// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package service
package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/beego/beego/v2/client/orm"

	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
)

var pathSeparatorRegex = regexp.MustCompile(`[/\\]`)

type FileService interface {
	UploadFile(bucket, name string, content []byte) (string, error)
	DownloadFile(bucket, name string) ([]byte, error)
	DeleteFile(bucket, name string) error
	Exist(bucket string, name string) (bool, error)
}

var _ FileService = &FileServiceImpl{}

func NewFileService() *FileServiceImpl {
	return &FileServiceImpl{
		fd: dao.NewFileDao(),
	}
}

type FileServiceImpl struct {
	fd *dao.FileDao
}

func (f *FileServiceImpl) Exist(bucket string, name string) (bool, error) {
	return f.fd.Exist(bucket, name)
}

func (f *FileServiceImpl) DeleteFile(bucket, name string) error {
	cleanedFileName, err := f.cleanFileName(name)
	if err != nil {
		return err
	}
	file := &db.File{
		Bucket: bucket,
		Name:   cleanedFileName,
	}
	return f.fd.Delete(file, "Bucket", "Name")
}

// 清理和验证文件名，防止路径遍历攻击
func (f *FileServiceImpl) cleanFileName(fileName string) (string, error) {
	cleanedFileName := filepath.Clean(fileName)
	if cleanedFileName == "" {
		logger.Errorf("invalid file name: %s", fileName)
		return "", errors.New("invalid file name")
	}

	if pathSeparatorRegex.MatchString(cleanedFileName) {
		logger.Errorf("invalid file name: %s", fileName)
		return "", errors.New("invalid file name")
	}

	return cleanedFileName, nil
}

func (f *FileServiceImpl) insertOrUpdate(new *db.File) error {
	old := &db.File{
		Bucket: new.Bucket,
		Name:   new.Name,
	}
	err := f.fd.Get(old, "Bucket", "Name")
	if err != nil && err != orm.ErrNoRows {
		return err
	}
	if err == orm.ErrNoRows {
		new.CreatedAt = time.Now().Format(time.DateTime)
		return f.fd.Insert(new)
	}
	old.Content = new.Content
	old.Size = new.Size
	return f.fd.Update(old)
}

func (f *FileServiceImpl) UploadFile(bucket, name string, content []byte) (string, error) {
	cleanedFileName, err := f.cleanFileName(name)
	if err != nil {
		return "", err
	}
	file := &db.File{
		Bucket:  bucket,
		Name:    cleanedFileName,
		Content: content,
		Size:    int64(len(content)),
	}

	if err := f.insertOrUpdate(file); err != nil {
		logger.Errorf("upload file to OSS error: %v", err)
		return "", errors.New("files upload to OSS failed")
	}
	return fmt.Sprintf("/%s/%s", bucket, cleanedFileName), nil
}

func (f *FileServiceImpl) DownloadFile(bucket, name string) ([]byte, error) {
	cleanedFileName, err := f.cleanFileName(name)
	if err != nil {
		return nil, err
	}
	file := &db.File{
		Bucket: bucket,
		Name:   cleanedFileName,
	}
	if err := f.fd.Get(file, "Bucket", "Name"); err != nil {
		logger.Errorf("download file from OSS failed: %v", err)
		return nil, errors.New("failed to download file from OSS")
	}
	return file.Content, nil
}
