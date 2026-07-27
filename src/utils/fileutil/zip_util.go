/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package fileutil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateZipFile(targetZipPath, sourceDir string) error {
	zipFile, err := os.OpenFile(targetZipPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, PermissionForLogFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", sourceDir, err)
	}

	for _, file := range files {
		// 只处理文件，跳过目录
		if !file.IsDir() {
			filePath := filepath.Join(sourceDir, file.Name())
			if targetZipPath == filePath {
				continue
			}
			if err = addFileToZip(zipWriter, filePath, file.Name()); err != nil {
				return fmt.Errorf("failed to add %s to ZIP: %v", file.Name(), err)
			}
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipFile, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, file)
	return err
}
