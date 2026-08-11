// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package service
package service

import (
	"sort"
	"sync"
	"time"
)

const (
	authCacheTTL        = 30 * time.Minute
	authCacheCapacity   = 1000
	authCacheEvictCount = 500
)

// cacheEntry 鉴权缓存条目
type cacheEntry struct {
	result   bool
	expireAt time.Time
}

// authCache 进程内鉴权缓存，RWMutex+map，读 RLock、写/清理 Lock
type authCache struct {
	sync.RWMutex
	entries map[string]cacheEntry
}

var (
	authCacheInstance *authCache
	authCacheOnce     sync.Once
	// authImportLock 导入管理 Service 持写锁、鉴权回源段持读锁，杜绝清表窗口期逃生态误放行
	authImportLock sync.RWMutex
)

// getAuthCache 获取鉴权缓存单例
func getAuthCache() *authCache {
	authCacheOnce.Do(func() {
		authCacheInstance = &authCache{
			entries: make(map[string]cacheEntry),
		}
	})
	return authCacheInstance
}

// get 查询缓存，过期视为 miss
func (c *authCache) get(key string) (bool, bool) {
	c.RLock()
	defer c.RUnlock()
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expireAt) {
		return false, false
	}
	return entry.result, true
}

// set 写入缓存，写入后容量超限按 expireAt 升序惰性清理最旧 500 条，清理与写入同一 Lock 内完成
func (c *authCache) set(key string, result bool) {
	c.Lock()
	defer c.Unlock()
	c.entries[key] = cacheEntry{result: result, expireAt: time.Now().Add(authCacheTTL)}
	if len(c.entries) <= authCacheCapacity {
		return
	}
	type kv struct {
		key      string
		expireAt time.Time
	}
	items := make([]kv, 0, len(c.entries))
	for k, v := range c.entries {
		items = append(items, kv{key: k, expireAt: v.expireAt})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].expireAt.Before(items[j].expireAt)
	})
	for i := 0; i < authCacheEvictCount && i < len(items); i++ {
		delete(c.entries, items[i].key)
	}
}
