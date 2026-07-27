/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package event

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"code.huawei.com/fusionstage/auditlog"

	"GIDS/common/constants"
	"GIDS/common/logger"
	"GIDS/models/events"
	util "GIDS/utils/fileutil"
)

const (
	FileMaxSize   = 20 * 1024 * 1024
	FileMaxNum    = 5
	FileRemainDay = 90
	hoursPerDay   = 24
)

// const
const (
	NanoFirstThreeDigit  = 3
	ZipFileDefaultLength = 27
	YearMonthDayStartPos = 6
	YearMonthDayEndPos   = 14
)

type localEventStorage struct {
	engine                auditlog.Logger
	name                  string
	logFile               string
	logFileHandle         *os.File
	routineTriggerDeleter []FileDeleter
	recordTriggerDeleter  []FileDeleter
}

// NewLocalEventStorage 创建LocalEventStorage
func NewLocalEventStorage(name string, logFile string) *localEventStorage {
	storage := &localEventStorage{
		name:    name,
		logFile: logFile,
		engine:  auditlog.NewLoggerBase(constants.ComponentName),
		routineTriggerDeleter: []FileDeleter{
			NewDayBasedDeleter(FileRemainDay),
			NewCountBasedDeleter(FileMaxNum),
		},
		recordTriggerDeleter: []FileDeleter{
			NewCountBasedDeleter(FileMaxNum),
		},
	}

	if logFile == "" {
		sink := auditlog.NewWriterSink(os.Stdout)
		storage.engine.RegisterSink(sink)
		return storage
	}

	if err := storage.createFileWhenNotExist(); err != nil {
		logger.Errorf("[localEventStorage] createFileWhenNotExist failed, filePath: [%s], err: [%v]", storage.logFile, err)
		sink := auditlog.NewWriterSink(os.Stdout)
		storage.engine.RegisterSink(sink)
		return storage
	}

	go func() {
		checkDelTicker := time.Tick(1 * time.Hour)

		for range checkDelTicker {
			logger.Infof("will exec delete old back event file. CurrentTime: %s", time.Now().String())
			storage.deleteOldBackFiles(true)
		}
	}()
	return storage
}

func (storage *localEventStorage) Record(event *events.Info) error {
	if storage.needRollOver() {
		storage.rollOver()
	}

	if _, err := storage.engine.Print(event.ToJSON()); err != nil {
		logger.Errorf("[localEventStorage] record content[%s] failed, err: %v", event.ToJSON(), err)
		return err
	}

	return nil
}

// implement io.Writer interface
func (storage *localEventStorage) Write(p []byte) (int, error) {
	if storage.needRollOver() {
		storage.rollOver()
	}

	return storage.engine.Print(p)
}

func (storage *localEventStorage) createFileWhenNotExist() error {
	_, err := os.Stat(storage.logFile)
	if err != nil && !os.IsNotExist(err) {
		logger.Errorf("[localEventStorage] failed to stat[%s], err: [%v]", storage.logFile, err)
		return err
	}

	if os.IsNotExist(err) {
		dir := filepath.Dir(storage.logFile)
		if err := os.MkdirAll(dir, util.PermissionForLogDir); err != nil {
			logger.Errorf("[localEventStorage] failed to createDirectory[%s], err: [%v]", dir, err)
			return err
		}
	}

	file, err := os.OpenFile(storage.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, util.PermissionForLogFile)
	if err != nil {
		logger.Errorf("[localEventStorage] failed to createFile[%s], err: [%v]", storage.logFile, err)
		return err
	}
	defer file.Close()
	storage.logFileHandle = file
	sink := auditlog.NewWriterSink(file)
	storage.engine.RegisterSink(sink)
	return nil
}
func (storage *localEventStorage) Close() {
	if storage.logFileHandle != nil {
		if err := storage.logFileHandle.Close(); err != nil {
			logger.Errorf("[localEventStorage] failed to close log file: %v", err)
		}
	}
}
func (storage *localEventStorage) needRollOver() bool {
	fileInfo, err := os.Stat(storage.logFile)
	if err != nil {
		logger.Errorf("[localEventStorage] stat file %s error: %v, skip rollOver localEventStorage", storage.logFile, err)
		return false
	}

	if fileInfo.Size() > FileMaxSize {
		return true
	}

	return false
}

func (storage *localEventStorage) rollOver() {
	dir := filepath.Dir(storage.logFile)

	now := time.Now()
	ymdhms := now.Format("20060102150405")
	ms := strconv.Itoa(now.Nanosecond())[:NanoFirstThreeDigit]

	backFile := dir + string(filepath.Separator) + "event_" + ymdhms + ms + ".zip"
	logTempFile := dir + string(filepath.Separator) + "event_" + ymdhms + ms + ".log"

	// 1. 复制需要转储的Event文件为临时文件
	if _, err := util.CopyFile(logTempFile, storage.logFile); err != nil {
		logger.Errorf("[localEventStorage] copy file %s to temp error: %v, skip rollOver localEventStorage", storage.logFile, err)
		return
	}

	// 2. 清空Event文件
	if err := os.Truncate(storage.logFile, 0); err != nil {
		logger.Errorf("[localEventStorage]truncate config.LogFile %s error: %v", storage.logFile, err)
		return
	}

	// 3. 压缩临时文件为最终的转储文件
	if err := util.ZipFile(logTempFile, backFile); err != nil {
		logger.Errorf("[localEventStorage] zip file %s to %s error: %v", storage.logFile, backFile, err)

		if err = os.Remove(logTempFile); err != nil {
			logger.Errorf("[localEventStorage] remove logTempFile %s error: %v", logTempFile, err)
		}
		return
	}

	if err := os.Chmod(backFile, util.PermissionForZipFile); err != nil {
		logger.Errorf("[localEventStorage] Chmod file[%s] failed, err: %v", backFile, err)
		return
	}

	// 4. 移除临时文件
	if err := os.Remove(logTempFile); err != nil {
		logger.Errorf("[localEventStorage] remove logTempFile %s error: %v", logTempFile, err)
		return
	}

	storage.deleteOldBackFiles(true)
}

func (storage *localEventStorage) deleteOldBackFiles(isRemainDayCheck bool) {
	dir := filepath.Dir(storage.logFile)
	backFiles := getBackFiles(dir)

	if len(backFiles) == 0 {
		logger.Infof("[localEventStorage] no zip Log to delete, logFile: %s", storage.logFile)
		return
	}

	sort.Strings(backFiles)

	deleters := storage.recordTriggerDeleter
	if isRemainDayCheck {
		deleters = storage.routineTriggerDeleter
	}

	for _, deleter := range deleters {
		backFiles = deleter.Handle(backFiles, dir)
	}
}

func getBackFiles(dir string) []string {
	var backFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Errorf("[localEventStorage] List dir[%s] error: %v", dir, err)
		return backFiles
	}

	for _, file := range files {
		match, err := regexp.Match("event_[\\d]{17}\\.zip", []byte(file.Name()))
		if err != nil {
			logger.Errorf("[localEventStorage] regexp fileName[%s] error!, err: %v", file.Name(), err)
			continue
		}
		if match {
			backFiles = append(backFiles, file.Name())
		}
	}

	return backFiles
}

type FileDeleter interface {
	Handle(backFiles []string, dir string) []string
}

type DayBasedDeleter struct {
	remainDay int
}

func NewDayBasedDeleter(remainDay int) *DayBasedDeleter {
	return &DayBasedDeleter{remainDay: remainDay}
}

func (d *DayBasedDeleter) Handle(backFiles []string, dir string) []string {
	duration, err := time.ParseDuration(fmt.Sprintf("-%dh", FileRemainDay*hoursPerDay))
	if err != nil {
		logger.Errorf("[localEventStorage] ParseDuration parses a duration string occur error: %v", err)
		return backFiles
	}

	oldestZipFileTimeExpected, err := strconv.Atoi(time.Now().Add(duration).Format("20060102"))
	if err != nil {
		logger.Errorf("[localEventStorage] convert nowTime to Expected oldestZipFileTime error: %v", err)
		return backFiles
	}

	var newBackFiles []string
	for _, fileName := range backFiles {
		length := len(fileName)
		if length < ZipFileDefaultLength {
			logger.Errorf("[localEventStorage] fileName[%s] invalid", fileName)
			continue
		}

		// 取zip文件中的年月日
		zipFileTimeByInt, err := strconv.Atoi(fileName[YearMonthDayStartPos:YearMonthDayEndPos])
		if err != nil {
			logger.Errorf("[localEventStorage] convert zipFileTime error: %v", err)
			continue
		}

		if zipFileTimeByInt-oldestZipFileTimeExpected <= 0 {
			logger.Infof("[localEventStorage] exec delete oldFileName: %s", fileName)
			if err := os.Remove(dir + string(filepath.Separator) + fileName); err != nil {
				logger.Errorf("[localEventStorage] remove old file error: %v", err)
				continue
			}
		}

		newBackFiles = append(newBackFiles, fileName)
	}

	return newBackFiles
}

type CountBasedDeleter struct {
	maxNum int
}

func NewCountBasedDeleter(maxNum int) *CountBasedDeleter {
	return &CountBasedDeleter{maxNum: maxNum}
}

func (c *CountBasedDeleter) Handle(backFiles []string, dir string) []string {
	if len(backFiles) < c.maxNum {
		return backFiles
	}

	shouldRemainFrom := len(backFiles) - c.maxNum + 1
	for i := 0; i < shouldRemainFrom; i++ {
		if err := os.Remove(dir + string(filepath.Separator) + backFiles[i]); err != nil {
			logger.Errorf("[localEventStorage] remove old file[%s] error: %v", backFiles[i], err)
			continue
		}
	}

	return backFiles[shouldRemainFrom:]
}
