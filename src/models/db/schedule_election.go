/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package db

import (
	"encoding/json"
	"fmt"
)

// ScheduleElection 定时任务选主
type ScheduleElection struct {
	Ip  string `json:"ip"`
	Mac string `json:"mac"`
	Id  int    `json:"Id"`
}

func (s *ScheduleElection) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ScheduleElection) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *ScheduleElection) GetKey() string {
	return fmt.Sprintf("gids.timerElection")
}
