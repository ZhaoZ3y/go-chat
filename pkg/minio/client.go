package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

// Config 结构体保持不变
type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

// Client 结构体保持不变
type Client struct {
	client     *minio.Client
	bucketName string
}

// NewClient 函数保持不变
func NewClient(config Config) (*Client, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return &Client{
		client:     client,
		bucketName: config.BucketName,
	}, nil
}

// UploadFile 上传文件，并返回 UploadInfo 以便获取 ETag
func (c *Client) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (minio.UploadInfo, error) {
	info, err := c.client.PutObject(ctx, c.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("failed to upload file: %w", err)
	}
	return info, nil
}

// DownloadFile 下载文件，直接返回文件内容的 ReadCloser
func (c *Client) DownloadFile(ctx context.Context, objectName string) (*minio.Object, error) {
	object, err := c.client.GetObject(ctx, c.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return object, nil
}

// DeleteFile 删除文件
func (c *Client) DeleteFile(ctx context.Context, objectName string) error {
	return c.client.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetFileInfo 获取文件信息
func (c *Client) GetFileInfo(ctx context.Context, objectName string) (minio.ObjectInfo, error) {
	return c.client.StatObject(ctx, c.bucketName, objectName, minio.StatObjectOptions{})
}

// GeneratePresignedURL 生成预签名 URL (保留此函数，可能在其他地方有用)
func (c *Client) GeneratePresignedURL(ctx context.Context, objectName string, expires time.Duration) (*url.URL, error) {
	return c.client.PresignedGetObject(ctx, c.bucketName, objectName, expires, nil)
}

// SetBucketLifecycle 设置存储桶的生命周期策略
func (c *Client) SetBucketLifecycle(ctx context.Context, days int) error {
	cfg := lifecycle.NewConfiguration()
	cfg.Rules = []lifecycle.Rule{
		{
			ID:     fmt.Sprintf("expire-after-%d-days", days),
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: lifecycle.ExpirationDays(days),
			},
			RuleFilter: lifecycle.Filter{
				Prefix: "", // 应用于桶内所有对象
			},
		},
	}
	return c.client.SetBucketLifecycle(ctx, c.bucketName, cfg)
}
