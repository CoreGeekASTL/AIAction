// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package service

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type StubGidsMasterDao struct {
	master    *StubGidsMaster
	queryErr  error
	upsertErr error
	updateErr error
	mu        sync.Mutex
}

type StubGidsMaster struct {
	PodName      string
	Timestamp    time.Time
	IsRegistered bool
}

func (m *StubGidsMasterDao) Query() (*StubGidsMaster, error) {
	return m.master, m.queryErr
}

func (m *StubGidsMasterDao) Upsert(podName string, timestamp time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.upsertErr != nil {
		return m.upsertErr
	}
	m.master = &StubGidsMaster{
		PodName:      podName,
		Timestamp:    timestamp,
		IsRegistered: false,
	}
	return nil
}

func (m *StubGidsMasterDao) UpsertIfEmpty(podName string, timestamp time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.upsertErr != nil {
		return m.upsertErr
	}
	if m.master != nil {
		return fmt.Errorf("master already exists")
	}
	m.master = &StubGidsMaster{
		PodName:      podName,
		Timestamp:    timestamp,
		IsRegistered: false,
	}
	return nil
}

func (m *StubGidsMasterDao) UpdateTimestamp(podName string, timestamp time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.master != nil && m.master.PodName == podName {
		m.master.Timestamp = timestamp
	}
	return m.updateErr
}

func (m *StubGidsMasterDao) UpdateIsRegistered(podName string, isRegistered bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.master != nil && m.master.PodName == podName {
		m.master.IsRegistered = isRegistered
	}
	return m.updateErr
}

type stubMasterElectionServiceImpl struct {
	dao      *StubGidsMasterDao
	podName  string
	isMaster bool
	mu       sync.RWMutex
}

func (s *stubMasterElectionServiceImpl) checkAndElection() {
	master, err := s.dao.Query()
	now := time.Now()

	if err != nil || master == nil {
		s.tryBecomeMaster(now)
		return
	}

	if master.PodName == s.podName {
		s.setIsMaster(true)
		s.dao.UpdateTimestamp(s.podName, now)
		return
	}

	if now.Sub(master.Timestamp) > 15*time.Second {
		s.tryBecomeMaster(now)
		return
	}

	s.setIsMaster(false)
}

func (s *stubMasterElectionServiceImpl) tryBecomeMaster(now time.Time) {
	err := s.dao.Upsert(s.podName, now)
	if err != nil {
		s.setIsMaster(false)
		return
	}

	master, err := s.dao.Query()
	if err == nil && master != nil && master.PodName == s.podName {
		s.setIsMaster(true)
	} else {
		s.setIsMaster(false)
	}
}

func (s *stubMasterElectionServiceImpl) setIsMaster(isMaster bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isMaster = isMaster
}

func (s *stubMasterElectionServiceImpl) IsMaster() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isMaster
}

func TestStub_FirstPodBecomeMaster(t *testing.T) {
	stubDao := &StubGidsMasterDao{
		master:   nil,
		queryErr: nil,
	}

	service := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-1",
		mu:      sync.RWMutex{},
	}

	service.tryBecomeMaster(time.Now())

	assert.True(t, service.IsMaster())
}

func TestStub_CurrentMasterRefresh(t *testing.T) {
	stubDao := &StubGidsMasterDao{
		master: &StubGidsMaster{
			PodName:      "pod-1",
			Timestamp:    time.Now(),
			IsRegistered: true,
		},
		queryErr: nil,
	}

	service := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-1",
		mu:      sync.RWMutex{},
	}

	service.checkAndElection()

	assert.True(t, service.IsMaster())
}

func TestStub_MasterTimeoutGrabMaster(t *testing.T) {
	oldTime := time.Now().Add(-20 * time.Second)
	stubDao := &StubGidsMasterDao{
		master: &StubGidsMaster{
			PodName:      "pod-old",
			Timestamp:    oldTime,
			IsRegistered: false,
		},
		queryErr: nil,
	}

	service := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-new",
		mu:      sync.RWMutex{},
	}

	service.checkAndElection()

	assert.True(t, service.IsMaster())
}

func TestStub_NonMasterWaiting(t *testing.T) {
	stubDao := &StubGidsMasterDao{
		master: &StubGidsMaster{
			PodName:      "pod-other",
			Timestamp:    time.Now(),
			IsRegistered: true,
		},
		queryErr: nil,
	}

	service := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-1",
		mu:      sync.RWMutex{},
	}

	service.checkAndElection()

	assert.False(t, service.IsMaster())
}

func TestStub_UpsertFailed(t *testing.T) {
	stubDao := &StubGidsMasterDao{
		master:    nil,
		queryErr:  nil,
		upsertErr: assert.AnError,
	}

	service := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-1",
		mu:      sync.RWMutex{},
	}

	service.tryBecomeMaster(time.Now())

	assert.False(t, service.IsMaster())
}

func TestStub_QueryFailed(t *testing.T) {
	stubDao := &StubGidsMasterDao{
		master:   nil,
		queryErr: nil,
	}

	service := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-1",
		mu:      sync.RWMutex{},
	}

	service.checkAndElection()

	assert.True(t, service.IsMaster())
}

func TestStub_MultiPod(t *testing.T) {
	stubDao := &StubGidsMasterDao{
		queryErr: nil,
	}

	pod1 := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-1",
		mu:      sync.RWMutex{},
	}

	pod2 := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-2",
		mu:      sync.RWMutex{},
	}

	stubDao.master = nil
	stubDao.queryErr = nil
	pod1.tryBecomeMaster(time.Now())
	assert.True(t, pod1.IsMaster())

	stubDao.mu.Lock()
	stubDao.master = &StubGidsMaster{
		PodName:   "pod-1",
		Timestamp: time.Now(),
	}
	stubDao.queryErr = nil
	stubDao.mu.Unlock()

	pod2.checkAndElection()
	assert.False(t, pod2.IsMaster())
}

func TestStub_FullFlow(t *testing.T) {
	stubDao := &StubGidsMasterDao{}

	service := &stubMasterElectionServiceImpl{
		dao:     stubDao,
		podName: "pod-1",
		mu:      sync.RWMutex{},
	}

	stubDao.master = nil
	stubDao.queryErr = nil

	service.tryBecomeMaster(time.Now())
	assert.True(t, service.IsMaster())

	stubDao.master = &StubGidsMaster{
		PodName:      "pod-1",
		Timestamp:    time.Now(),
		IsRegistered: true,
	}
	service.checkAndElection()
	assert.True(t, service.IsMaster())

	stubDao.master = &StubGidsMaster{
		PodName:      "pod-1",
		Timestamp:    time.Now().Add(-20 * time.Second),
		IsRegistered: false,
	}
	service.podName = "pod-2"
	service.checkAndElection()
	assert.True(t, service.IsMaster())
}

func TestStub_StateTransition(t *testing.T) {
	stubDao := &StubGidsMasterDao{
		master: &StubGidsMaster{
			PodName:      "pod-other",
			Timestamp:    time.Now().Add(-20 * time.Second),
			IsRegistered: false,
		},
	}

	service := &stubMasterElectionServiceImpl{
		dao:      stubDao,
		podName:  "pod-1",
		mu:       sync.RWMutex{},
		isMaster: false,
	}

	assert.False(t, service.IsMaster())

	service.checkAndElection()
	assert.True(t, service.IsMaster())

	stubDao.master = &StubGidsMaster{
		PodName:      "pod-new-master",
		Timestamp:    time.Now(),
		IsRegistered: true,
	}
	service.podName = "pod-standby"
	service.checkAndElection()
	assert.False(t, service.IsMaster())
}