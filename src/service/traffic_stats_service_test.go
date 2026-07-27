/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2026. All rights reserved.
 */

package service

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"GIDS/dao"
)

/*
* 测试用例描述： TestTrafficStatsServiceImplGetOnline
* 预置条件：无
* 操作步骤：
*     1. 测试按照sqlConfig的sql语句配置 去获取指标值
* 预期结果：
*     1. 能正确执行配置的sql语句 并返回结果
* 修改历史：
*     1. 2025-12-14 新建测试用例
 */
func TestTrafficStatsServiceImplGetOnline(t *testing.T) {

	type args struct {
		startTime string
		endTime   string
	}
	tests := []struct {
		name      string
		sd        *dao.SessionStatsDao
		sqlConfig *SQLConfig
		args      args
		want      []Res
	}{
		{
			name: "getTrafficObjAndValue",
			sd: &dao.SessionStatsDao{
				BaseInterface: &dao.DoNothingBase{},
			},
			sqlConfig: &SQLConfig{
				Queries: map[string]struct {
					SQL    string   `yaml:"sql"`
					Params []string `yaml:"params"`
				}{
					getOnline: {
						SQL:    "SELECT * FROM table1 WHERE id = ?",
						Params: []string{"id"},
					},
					getOnlineOfModel: {
						SQL:    "INSERT INTO table2 (name, age) VALUES (?, ?)",
						Params: []string{"name", "age"},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TrafficStatsServiceImpl{
				sd:        tt.sd,
				sqlConfig: tt.sqlConfig,
			}
			assert.Equalf(t, tt.want, s.GetOnline(tt.args.startTime, tt.args.endTime), "GetOnline(%v, %v)", tt.args.startTime, tt.args.endTime)
			assert.Equalf(t, tt.want, s.GetOnlineOfModel(tt.args.startTime, tt.args.endTime), "GetOnlineOfModel(%v, %v)", tt.args.startTime, tt.args.endTime)
			assert.Equalf(t, tt.want, s.GetTrafficOfApp(tt.args.startTime, tt.args.endTime), "GetTrafficOfApp(%v, %v)", tt.args.startTime, tt.args.endTime)
		})
	}
}
