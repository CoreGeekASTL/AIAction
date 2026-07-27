// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package fileutil 提供文件操作和压缩工具函数。
package fileutil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"runtime"

	"code.huawei.com/fusionstage/auditlog"

	"GIDS/common/logger"
)

const (
	PermissionForLogFile = 0600
	PermissionForZipFile = 0400
	PermissionForLogDir  = 0755
)

func CopyFile(dstFile string, srcFile string) (int64, error) {
	srcFileName, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := srcFileName.Close(); err != nil {
			logger.Errorf("fail to close fileHandler, errInfo: %v", err)
		}
	}()

	dstFileName, err := os.OpenFile(dstFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, PermissionForLogFile)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := dstFileName.Close(); err != nil {
			logger.Errorf("fail to close fileHandler, errInfo: %v", err)
		}
	}()

	return io.Copy(dstFileName, srcFileName)
}

func ZipFile(source, target string) error {
	targetFile, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, auditlog.PERMISSION0440)
	if err != nil {
		return err
	}
	defer func() {
		if err := targetFile.Close(); err != nil {
			logger.Errorf("fail to close fileHandler, errInfo: %v", err)
		}
	}()
	archiveFile := zip.NewWriter(targetFile)
	defer func() {
		if err := archiveFile.Close(); err != nil {
			logger.Errorf("fail to close fileHandler, errInfo: %v", err)
		}
	}()
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			logger.Errorf("fail to close fileHandler, errInfo: %v", err)
		}
	}()

	// get the file information
	info, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// Deflate: 压缩进行转储
	header.Method = zip.Deflate

	writer, err := archiveFile.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, sourceFile)
	if err != nil {
		return err
	}
	err = copyFileOwnInfo(source, target)

	return err
}

func CloseFileHandle(file io.Closer) {
	if file == nil {
		logger.Errorf("fileHandler is nil.")
		return
	}
	if err := file.Close(); err != nil {
		logger.Errorf("fail to close fileHandler, errInfo: %v", err)
	}
}

func copyFileOwnInfo(src, dest string) error {
	uidSrc, gidSrc, err := getFileOwnInfo(src)
	if err != nil {
		return err
	}
	return setFileOwnInfo(dest, int(uidSrc), int(gidSrc))
}

type owner struct {
	uid uint32
	gid uint32
}

func getFileOwnInfo(file string) (uid, gid uint32, err error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return auditlog.DefaultUserID, auditlog.DefaultGroupID, err
	}

	o, ok := getOwner(fileInfo)
	if !ok {
		return auditlog.DefaultUserID, auditlog.DefaultGroupID, fmt.Errorf("failed to getFileOwnInfo")
	}
	return o.uid, o.gid, nil
}

func getOwner(_ os.FileInfo) (owner, bool) {
	return owner{
		uid: auditlog.DefaultUserID,
		gid: auditlog.DefaultUserID,
	}, true
}

func setFileOwnInfo(file string, uid, gid int) error {
	if runtime.GOOS == "windows" {
		return nil
	}

	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		if err := fp.Close(); err != nil {
			logger.Errorf("fail to close fileHandler, errInfo: %v", err)
		}
	}()
	return fp.Chown(uid, gid)
}
