// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

package service

import (
	"sort"
	"sync"
	"time"
)

const (
	// authCacheCapacity 鉴权缓存容量上限，超限触发惰性清理
	authCacheCapacity = 1000
	// authCacheCleanCount 每次清理最旧条目数
	authCacheCleanCount = 500
	// authCacheTTL 缓存条目存活时长，过期回源 DB
	authCacheTTL = 30 * time.Minute
)

// cacheEntry 缓存项：鉴权结果（命中/未命中/逃生态放行标记）+ 过期时间
type cacheEntry struct {
	isAllowed bool
	expireAt  time.Time
}

// authCache 鉴权结果内存缓存，读写锁保护，无独立 goroutine，清理在写锁内完成
type authCache struct {
	sync.RWMutex
	items map[string]cacheEntry
}

func newAuthCache() *authCache {
	return &authCache{items: make(map[string]cacheEntry)}
}

// get 查询缓存，命中且未过期返回结果与 true
func (c *authCache) get(key string) (bool, bool) {
	c.RLock()
	defer c.RUnlock()
	entry, ok := c.items[key]
	if !ok || time.Now().After(entry.expireAt) {
		return false, false
	}
	return entry.isAllowed, true
}

// set 写入缓存，写后容量超限则按过期时间升序清理最旧条目
func (c *authCache) set(key string, isAllowed bool) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = cacheEntry{isAllowed: isAllowed, expireAt: time.Now().Add(authCacheTTL)}
	if len(c.items) > authCacheCapacity {
		c.cleanLocked()
	}
}

// clear 清空缓存，白名单导入成功后调用，保证立即生效
func (c *authCache) clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[string]cacheEntry)
}

// cleanLocked 惰性清理：按 expireAt 升序删除最旧 authCacheCleanCount 条，调用方须持有写锁
func (c *authCache) cleanLocked() {
	type expiringItem struct {
		key      string
		expireAt time.Time
	}
	entries := make([]expiringItem, 0, len(c.items))
	for key, entry := range c.items {
		entries = append(entries, expiringItem{key: key, expireAt: entry.expireAt})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].expireAt.Before(entries[j].expireAt) })
	for i := 0; i < authCacheCleanCount && i < len(entries); i++ {
		delete(c.items, entries[i].key)
	}
}
