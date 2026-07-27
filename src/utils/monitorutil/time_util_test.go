/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package monitorutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetLastFiveMinuteWindow(t *testing.T) {
	tests := []struct {
		name    string
		cur     time.Time
		startAt string
		endAt   string
	}{
		{
			name:    "curLastFiveMinute01",
			cur:     time.Date(2025, 12, 8, 16, 14, 18, 0, time.Local),
			startAt: "2025-12-08 16:05:00",
			endAt:   "2025-12-08 16:09:59",
		},
		{
			name:    "curLastFiveMinute02",
			cur:     time.Date(2025, 12, 8, 16, 03, 18, 0, time.Local),
			startAt: "2025-12-08 15:55:00",
			endAt:   "2025-12-08 15:59:59",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retStartAt, retEndAt := GetLastFiveMinuteWindow(&tt.cur)
			assert.Equal(t, tt.startAt, retStartAt)
			assert.Equal(t, tt.endAt, retEndAt)
		})
	}
}
