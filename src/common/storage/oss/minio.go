// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package oss
package oss

import (
	"context"
	"errors"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"GIDS/common/conf"
	"GIDS/common/logger"
)

var client Client

type Client interface {
	PutObject(ctx context.Context, object PutObject) error
	EnsureBucket(ctx context.Context, bucket string) error
	IsOnline() bool
	DeleteObject(background context.Context, bucket string, filename string) error
	GetObject(ctx context.Context, options GetObjectOptions) (io.ReadCloser, error)
}

func Init(config conf.OSSConfig) error {
	c, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, config.Token),
		Secure: false,
	})
	if err != nil {
		logger.Errorf("init minio failed! err is %v", err)
		return err
	}
	client = &ossClient{client: c}
	if !client.IsOnline() {
		err = errors.New("minio healthcheck failed")
		logger.Errorf("%v", err)
		return err
	}
	return nil
}

func Instance() Client {
	return client
}

type ossClient struct {
	client *minio.Client
}

type PutObject struct {
	BucketName string
	FileName   string
	File       io.Reader
	Size       int64
}

type GetObjectOptions struct {
	BucketName string
	FileName   string
}

func (c *ossClient) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := c.client.BucketExists(ctx, bucket)
	if err != nil {
		logger.Errorf("check bucket exist error! err is %v", err)
		return err
	}
	if !exists {
		err := c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Errorf("make bucket error! err is %v", err)
			return err
		}
	}
	return nil
}

func (c *ossClient) IsOnline() bool {
	return c.client.IsOnline()
}

func (c *ossClient) PutObject(ctx context.Context, object PutObject) error {
	info, err := c.client.PutObject(ctx,
		object.BucketName,
		object.FileName,
		object.File,
		object.Size,
		minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	logger.Infof("put object success, return message is %v", info)
	return nil
}

func (c *ossClient) DeleteObject(background context.Context, bucketName string, filename string) error {
	err := c.client.RemoveObject(background, bucketName, filename, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		logger.Errorf("delete object error, err is %v", err)
		return err
	}
	return nil
}

func (c *ossClient) GetObject(ctx context.Context, options GetObjectOptions) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, options.BucketName, options.FileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
