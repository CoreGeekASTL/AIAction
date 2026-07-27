/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

// Package monitorutil 提供监控时间窗口计算工具。
package monitorutil

import "time"

const (
	fiveMinuteWindow        = 5
	fiveMinuteWindowEndMinute = 55
)

// GetLastFiveMinuteWindow 计算前一个整5分钟周期的开始时间 例如：19:34 -> [19:25:00, 19:30:00 ]
func GetLastFiveMinuteWindow(inputTime *time.Time) (string, string) {
	var now time.Time
	if inputTime != nil {
		now = *inputTime
	} else {
		now = time.Now()
	}
	minute := now.Minute()

	var startTime time.Time
	if minute < fiveMinuteWindow {
		startTime = time.Date(now.Year(), now.Month(), now.Day(),
			now.Hour()-1, fiveMinuteWindowEndMinute, 0, 0,
			now.Location(),
		)
	} else {
		startMinute := (minute/fiveMinuteWindow)*fiveMinuteWindow - fiveMinuteWindow
		startTime = time.Date(now.Year(), now.Month(), now.Day(),
			now.Hour(), startMinute, 0, 0,
			now.Location(),
		)
	}

	endTime := startTime.Add(fiveMinuteWindow * time.Minute).Add(-time.Second)
	return startTime.Format(time.DateTime), endTime.Format(time.DateTime)
}
