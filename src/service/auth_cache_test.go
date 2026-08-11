// Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.

// Package service
package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
* 测试用例描述：TestAuthCacheGetSet
* 预置条件：新建缓存实例
* 操作步骤：
*     1. 查询不存在的键
*     2. 写入后查询
*     3. 写入过期条目后查询
* 预期结果：
*     1. miss
*     2. 命中并返回写入结果
*     3. 过期视为 miss
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthCacheGetSet(t *testing.T) {
	c := &authCache{entries: make(map[string]cacheEntry)}

	_, hit := c.get("missing")
	assert.False(t, hit)

	c.set("k1", true)
	result, hit := c.get("k1")
	assert.True(t, hit)
	assert.True(t, result)

	c.set("k2", false)
	result, hit = c.get("k2")
	assert.True(t, hit)
	assert.False(t, result)

	c.entries["k3"] = cacheEntry{result: true, expireAt: time.Now().Add(-time.Minute)}
	_, hit = c.get("k3")
	assert.False(t, hit)
}

/*
* 测试用例描述：TestAuthCacheEvict
* 预置条件：新建缓存实例
* 操作步骤：
*     1. 写入 capacity+1 条记录
* 预期结果：
*     1. 惰性清理最旧 500 条，容量降到 501
* 修改历史：
*     1. 2026-8-11 新建测试用例
 */
func TestAuthCacheEvict(t *testing.T) {
	c := &authCache{entries: make(map[string]cacheEntry)}
	for i := 0; i < authCacheCapacity; i++ {
		c.set(fmt.Sprintf("key-%04d", i), true)
	}
	assert.Equal(t, authCacheCapacity, len(c.entries))

	c.set("overflow", true)
	assert.Equal(t, authCacheCapacity-authCacheEvictCount+1, len(c.entries))

	_, hit := c.get("overflow")
	assert.True(t, hit)
}
