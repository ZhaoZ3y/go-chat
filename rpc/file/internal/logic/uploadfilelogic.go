package logic

import (
	pkgminio "IM/pkg/minio"
	"IM/pkg/model"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 服务端直接上传文件
func (l *UploadFileLogic) UploadFile(in *file.UploadFileRequest) (*file.UploadFileResponse, error) {
	// 1. 计算文件哈希
	hash := fmt.Sprintf("%x", md5.Sum(in.FileData))

	// 2. 检查文件是否已存在 (秒传逻辑)
	// 注意：这里的秒传是基于用户ID的，如果需要全局秒传，移除 UserID 条件
	var existingFile model.Files
	if err := l.svcCtx.DB.Where("hash = ? AND user_id = ?", hash, in.UserId).First(&existingFile).Error; err == nil {
		return &file.UploadFileResponse{
			FileId:  existingFile.Id,
			FileUrl: existingFile.FileUrl, // 数据库中已存正确的 MinIO URL
			Success: true,
			Message: "文件已存在，秒传成功",
		}, nil
	}

	// 3. 生成 MinIO 对象名称
	// in.Filename 是客户端的原始文件名, in.FileType 是文件分类 (如 "image", "video")
	objectName := pkgminio.GenerateObjectName(in.UserId, in.FileType, in.Filename)

	// 4. 上传文件到 MinIO
	reader := bytes.NewReader(in.FileData)
	err := l.svcCtx.Minio.UploadFile(l.ctx, objectName, reader, int64(len(in.FileData)), in.MimeType)
	if err != nil {
		l.Logger.Errorf("上传文件到 MinIO 失败: %v, 对象名: %s", err, objectName)
		return &file.UploadFileResponse{
			Success: false,
			Message: "上传文件到存储服务失败: " + err.Error(),
		}, nil
	}

	// 5. 生成文件访问 URL
	fileURL := l.buildMinioFileUrl(objectName)

	// 6. 保存文件记录到数据库
	fileRecord := model.Files{
		Filename:     filepath.Base(objectName), // 存储 MinIO 中的实际文件名 (通常是 time+random.ext)
		OriginalName: in.Filename,               // 存储用户上传的原始文件名
		FilePath:     objectName,                // 存储完整的 MinIO Object Key
		FileUrl:      fileURL,
		FileType:     in.FileType,
		FileSize:     int64(len(in.FileData)),
		MimeType:     in.MimeType,
		Hash:         hash,
		UserId:       in.UserId,
		Status:       1, // 1 表示正常
		CreateAt:     time.Now().Unix(),
		UpdateAt:     time.Now().Unix(),
	}

	if err := l.svcCtx.DB.Create(&fileRecord).Error; err != nil {
		l.Logger.Errorf("保存文件记录到数据库失败: %v, 文件信息: %+v", err, fileRecord)
		// 如果数据库保存失败，尝试删除已上传到 MinIO 的文件，尽力回滚
		if delErr := l.svcCtx.Minio.DeleteFile(l.ctx, objectName); delErr != nil {
			l.Logger.Errorf("数据库保存失败后，删除 MinIO 对象 %s 也失败: %v", objectName, delErr)
		}
		return &file.UploadFileResponse{
			Success: false,
			Message: "保存文件信息失败: " + err.Error(),
		}, nil
	}

	return &file.UploadFileResponse{
		FileId:  fileRecord.Id,
		FileUrl: fileRecord.FileUrl,
		Success: true,
		Message: "上传成功",
	}, nil
}

// 构建 MinIO 文件 URL 的辅助函数
func (l *UploadFileLogic) buildMinioFileUrl(objectName string) string {
	cfgMinio := l.svcCtx.Config.MinIO
	cfgStorage := l.svcCtx.Config.FileStorage

	// 优先使用 FileStorage.BaseURL (通常是CDN或反向代理地址)
	if cfgStorage.BaseURL != "" {
		return fmt.Sprintf("%s/%s/%s",
			strings.TrimSuffix(cfgStorage.BaseURL, "/"),
			cfgMinio.BucketName,
			objectName)
	}
	// 否则，直接使用 MinIO 端点
	scheme := "http"
	if cfgMinio.UseSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s",
		scheme,
		cfgMinio.Endpoint,
		cfgMinio.BucketName,
		objectName)
}
