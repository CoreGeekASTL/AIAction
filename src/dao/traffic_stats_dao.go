// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package dao

import (
	"fmt"

	"GIDS/models/db"
)

type MediaTrafficStatsDao struct {
	BaseInterface
}

func NewMediaTrafficStatsDao() *MediaTrafficStatsDao {
	dao := &MediaTrafficStatsDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.MediaTrafficStats{},
	}
	return dao
}

type ControlTrafficStatsDao struct {
	BaseInterface
}

func NewControlTrafficStatsDao() *ControlTrafficStatsDao {
	dao := &ControlTrafficStatsDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.ControlTrafficStats{},
	}
	return dao
}

type SessionStatsDao struct {
	BaseInterface
}

func NewSessionLogDao() *SessionStatsDao {
	dao := &SessionStatsDao{}
	dao.BaseInterface = &BaseDao{
		EntityType: &db.SessionStats{},
	}
	return dao
}

func (f *SessionStatsDao) Exist(tcpId string) (bool, error) {
	var exists bool
	err := f.QueryOne(&exists, "SELECT EXISTS(SELECT 1 FROM t_session_stats WHERE tcp_unique_id = ?)", tcpId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (f *SessionStatsDao) GetIdByTcpUniqueId(tcpUniqueId string) (int, error) {
	var id int
	err := f.QueryOne(&id, "SELECT id FROM t_session_stats WHERE tcp_unique_id = ?", tcpUniqueId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (f *SessionStatsDao) UpdatebySession(session *db.SessionStats) error {
	oldId, err := f.GetIdByTcpUniqueId(session.TcpUniqueId)
	if err != nil {
		return fmt.Errorf("get session failed: %w", err)
	}
	session.ID = oldId
	if err := f.Update(session, "finished_at"); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}
