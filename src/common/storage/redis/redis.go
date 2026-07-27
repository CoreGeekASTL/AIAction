// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package redis impl
package redis

import (
	"context"
	"encoding"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"GIDS/common/conf"
	"GIDS/common/logger"
	"GIDS/common/storage"
)

const maxRedisDB = 15

var client Client

// Init redis client
func Init(config conf.RedisConfig) error {
	if config.DB < 0 || config.DB > maxRedisDB {
		logger.Errorf("illegal redis db number[%s]", config.DB)
		return errors.New("illegal redis db number]")
	}

	client = New(&redis.Options{
		Addr:     config.Endpoint,
		Username: "",
		Password: "",
		DB:       config.DB,
	})
	err := client.Ping(context.Background())
	if err != nil {
		return err
	}
	logger.Infof("redis init success")
	return nil
}

func InitForTest(data map[string]interface{}) func() {
	return func() {
		client = nil
	}
}

// Instance return redis client
func Instance() Client {
	return client
}

// Object redis object
type Object interface {
	GetKey() string
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
type HFieldObject interface {
	Object
	GetField() string
}

// Client redis client
type Client interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, dst Object) error
	Exist(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
	Set(ctx context.Context, src Object) error
	SAdd(ctx context.Context, key string, value string) error
	SRem(ctx context.Context, key string, value string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SetIsMember(ctx context.Context, key string, value string) (bool, error)
	HSet(ctx context.Context, src HFieldObject) error
	HGet(ctx context.Context, dst HFieldObject) error
	HDel(ctx context.Context, src HFieldObject) error
	HKeyExists(ctx context.Context, key, field string) (bool, error)
	HGetAll(ctx context.Context, dst Object) error
	HKeys(ctx context.Context, key string) ([]string, error)
	SetNx(ctx context.Context, src Object, expiration time.Duration) (bool, error)
	SetWithExpiration(ctx context.Context, src Object, expiration time.Duration) error
}

type innerClient struct {
	client *redis.Client
}

func (c *innerClient) SetWithExpiration(ctx context.Context, src Object, expiration time.Duration) error {
	cmd := c.client.Set(ctx, src.GetKey(), src, expiration)
	return cmd.Err()
}

func (c *innerClient) SetNx(ctx context.Context, src Object, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, src.GetKey(), src, expiration).Result()
}

// New redis client
func New(option *redis.Options) Client {
	return &innerClient{client: redis.NewClient(
		option,
	)}
}

// Ping redis
func (c *innerClient) Ping(ctx context.Context) error {
	cmd := c.client.Ping(ctx)
	return cmd.Err()
}

// Get read a redis object
func (c *innerClient) Get(ctx context.Context, dst Object) error {
	result, err := c.client.Get(ctx, dst.GetKey()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return storage.ErrNotExist
		}
		return err
	}

	if err := dst.UnmarshalBinary([]byte(result)); err != nil {
		return err
	}

	return nil
}

// HGetAll read a redis object
func (c *innerClient) HGetAll(ctx context.Context, dst Object) error {
	cmd := c.client.HGetAll(ctx, dst.GetKey())
	err := cmd.Scan(dst)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return storage.ErrNotExist
		}
		return err
	}
	return nil
}

// Exist check redis object exist
func (c *innerClient) Exist(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result == 1, err
}

// Del delete redis object
func (c *innerClient) Del(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// Set write a redis object
func (c *innerClient) Set(ctx context.Context, src Object) error {
	cmd := c.client.Set(ctx, src.GetKey(), src, 0)
	return cmd.Err()
}

// SMembers query set members
func (c *innerClient) SMembers(ctx context.Context, key string) ([]string, error) {
	result, err := c.client.SMembers(ctx, key).Result()
	return result, wrapErr(err)
}

// HSet  set members
func (c *innerClient) HSet(ctx context.Context, src HFieldObject) error {
	_, err := c.client.HSet(ctx, src.GetKey(), src.GetField(), src).Result()
	return wrapErr(err)
}

// HKeyExists  get members
func (c *innerClient) HKeyExists(ctx context.Context, key, field string) (bool, error) {
	result, err := c.client.HExists(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, storage.ErrNotExist
		}
		return false, err
	}

	return result, nil
}

// HGet  get members
func (c *innerClient) HGet(ctx context.Context, dst HFieldObject) error {
	result, err := c.client.HGet(ctx, dst.GetKey(), dst.GetField()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return storage.ErrNotExist
		}
		return err
	}

	if err := dst.UnmarshalBinary([]byte(result)); err != nil {
		return err
	}

	return nil
}

func (c *innerClient) HDel(ctx context.Context, src HFieldObject) error {
	return c.client.HDel(ctx, src.GetKey(), src.GetField()).Err()
}

func (c *innerClient) HKeys(ctx context.Context, key string) ([]string, error) {
	fields, err := c.client.HKeys(ctx, key).Result()
	return fields, wrapErr(err)
}

func (c *innerClient) SAdd(ctx context.Context, key string, value string) error {
	return c.client.SAdd(ctx, key, value).Err()
}

func (c *innerClient) SRem(ctx context.Context, key string, value string) error {
	return c.client.SRem(ctx, key, value).Err()
}

func (c *innerClient) SetIsMember(ctx context.Context, key string, value string) (bool, error) {
	return c.client.SIsMember(ctx, key, value).Result()
}

func wrapErr(err error) error {
	if errors.Is(err, redis.Nil) {
		return storage.ErrNotExist
	}
	return err
}
