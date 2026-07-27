// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package scheduler
package scheduler

import (
	"sync"
	"time"

	"GIDS/common/constants"
	"GIDS/common/logger"
	"GIDS/service"
)

const hoursPerDay = 24

// DataCleanupScheduler 数据清理调度器
type DataCleanupScheduler struct {
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}

var globalScheduler *DataCleanupScheduler

func init() {
	globalScheduler = NewDataCleanupScheduler()
}

// NewDataCleanupScheduler 新建调度器实例
func NewDataCleanupScheduler() *DataCleanupScheduler {
	return &DataCleanupScheduler{
		stopChan: make(chan struct{}),
	}
}

func StartDataCleanupScheduler() {
	globalScheduler.Start()
}

func StopDataCleanupScheduler() {
	globalScheduler.Stop()
}

func (s *DataCleanupScheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		logger.Infof("Data cleanup scheduler already running")
		return
	}

	s.isRunning = true
	s.waitGroup.Add(1)
	go s.run()
	logger.Infof("Data cleanup scheduler started")
}

func (s *DataCleanupScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		logger.Infof("Data cleanup scheduler not running")
		return
	}

	close(s.stopChan)
	s.waitGroup.Wait()
	s.isRunning = false
	logger.Infof("Data cleanup scheduler stopped")
}

// 主运行逻辑
func (s *DataCleanupScheduler) run() {
	defer s.waitGroup.Done()

	const (
		maxRetries    = 3
		retryInterval = 10 * time.Minute
	)

	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		nextRun := s.calculateNextRunTime()
		sleepDuration := nextRun.Sub(time.Now())

		if sleepDuration < 0 {
			logger.Warnf("Next run time is in the past, resetting to next schedule cycle")
			continue
		}
		logger.Infof("Next cleanup scheduled at %v (in %v)", nextRun.Format(time.RFC3339), sleepDuration)

		timer.Reset(sleepDuration)
		select {
		case <-timer.C:
			s.mu.Lock()
			success := s.executeCleanup(maxRetries, retryInterval)
			s.mu.Unlock()
			if !success {
				logger.Errorf("Data cleanup failed, will retry at next schedule")
			}
		case <-s.stopChan:
			logger.Infof("Stopping scheduler...")
			return
		}
	}
}

// 计算下次执行时间
func (s *DataCleanupScheduler) calculateNextRunTime() time.Time {
	// 计算到下一个凌晨2点的时间
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())

	// 如果今天已经过了2点，则设置为明天2点
	if now.After(nextRun) {
		nextRun = nextRun.Add(hoursPerDay * time.Hour)
	}

	return nextRun
}

// 执行清理任务
func (s *DataCleanupScheduler) executeCleanup(maxRetries int, retryInterval time.Duration) bool {
	logger.Infof("Starting scheduled cleanup for %d months old data", constants.CleanupMonths)

	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		err := service.NewTrafficStatsService("").CleanOldStats(constants.CleanupMonths)
		if err == nil {
			logger.Infof("Cleanup completed successfully")
			return true
		}

		logger.Errorf("Cleanup failed (attempt %d/%d): %v", retryCount+1, maxRetries, err)
		if !s.sleepWithStopCheck(retryInterval) {
			// 如果在重试期间被停止，直接返回
			logger.Infof("Stopping scheduler during retry...")
			return false
		}
	}

	logger.Errorf("Max retries (%d) reached. Waiting for next schedule", maxRetries)
	return false
}

// 带停止检查的等待
func (s *DataCleanupScheduler) sleepWithStopCheck(duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-s.stopChan:
		logger.Infof("Stopping scheduler during retry...")
		return false
	}
}
