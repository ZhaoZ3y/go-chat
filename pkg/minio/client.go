package minio

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client     *minio.Client
	bucketName string
}

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

// NewMinIOClient 创建一个新的 MinIO 客户端
func NewMinIOClient(cfg Config) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("创建MinIO客户端失败: %v", err)
	}

	// 检查桶是否存在，不存在则创建
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("检查桶存在性失败: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("创建桶失败: %v", err)
		}
	}

	return &Client{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

// GetClient 获取原生minio客户端
func (m *Client) GetClient() *minio.Client {
	return m.client
}

// GetBucketName 获取桶名称
func (m *Client) GetBucketName() string {
	return m.bucketName
}

// UploadFile 上传文件到 MinIO
func (m *Client) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, objectSize, opts)
	if err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	return err
}

// DownloadFile 下载文件从 MinIO
func (m *Client) DownloadFile(ctx context.Context, objectName string) (*minio.Object, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}

	return object, nil
}

// DeleteFile 删除文件从 MinIO
func (m *Client) DeleteFile(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// GetPresignedURL 生成预签名的下载URL
func (m *Client) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}

	return url.String(), nil
}

// GetPresignedPutURL 生成预签名的上传URL
func (m *Client) GetPresignedPutURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	url, err := m.client.PresignedPutObject(ctx, m.bucketName, objectName, expires)
	if err != nil {
		return "", fmt.Errorf("生成预签名上传URL失败: %v", err)
	}

	return url.String(), nil
}

// FileExists 检查文件是否存在
func (m *Client) FileExists(ctx context.Context, objectName string) bool {
	_, err := m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
	return err == nil
}

// GenerateObjectName 生成对象名称
func GenerateObjectName(userID int64, fileType, originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s/%d/%d%s", fileType, userID, timestamp, ext)
}
