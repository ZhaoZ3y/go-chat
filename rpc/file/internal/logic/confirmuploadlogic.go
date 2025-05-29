package logic

import (
	"IM/pkg/model"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"path/filepath"
	"strings"
	"time"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmUploadLogic {
	return &ConfirmUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 确认上传 (配合预签名URL使用)
func (l *ConfirmUploadLogic) ConfirmUpload(in *file.ConfirmUploadRequest) (*file.ConfirmUploadResponse, error) {
	// 验证文件是否真的在 MinIO 中存在
	objInfo, err := l.svcCtx.Minio.GetClient().StatObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, in.FileKey, minio.StatObjectOptions{})
	if err != nil {
		l.Logger.Errorf("在 MinIO 中查询对象 %s 失败: %v", in.FileKey, err)
		return &file.ConfirmUploadResponse{
			Success: false,
			Message: "文件确认失败：无法在存储服务中找到对应文件或查询出错。",
		}, nil
	}

	fileSize := objInfo.Size
	mimeType := objInfo.ContentType
	if in.FileSize > 0 && in.FileSize != fileSize {
		l.Logger.Errorf("客户端提供的FileSize (%d) 与 MinIO 中的FileSize (%d) 不符，将使用 MinIO 的值。对象: %s", in.FileSize, fileSize, in.FileKey)
	}
	if in.MimeType != "" && in.MimeType != mimeType {
		l.Logger.Errorf("客户端提供的MimeType (%s) 与 MinIO 中的MimeType (%s) 不符，将使用 MinIO 的值。对象: %s", in.MimeType, mimeType, in.FileKey)
	}
	// 如果输入中没有，则使用MinIO的值
	if in.FileSize == 0 {
		in.FileSize = fileSize
	}
	if in.MimeType == "" {
		in.MimeType = mimeType
	}

	// 检查文件是否已通过哈希存在 (秒传逻辑)
	//  客户端在 ConfirmUpload 时应提供文件哈希。
	if in.Hash == "" {
		l.Logger.Errorf("ConfirmUpload 请求中缺少文件哈希，对象键: %s, 用户ID: %d", in.FileKey, in.UserId)
		return &file.ConfirmUploadResponse{Success: false, Message: "缺少文件哈希"}, nil
	}

	if in.Hash != "" {
		var existingFile model.Files
		if err := l.svcCtx.DB.Where("hash = ? AND user_id = ?", in.Hash, in.UserId).First(&existingFile).Error; err == nil {
			// 文件已存在，理论上客户端在 GeneratePresignedUrl 时就可能发现 (如果那时提供了hash)
			// 或者，这是并发上传了相同文件的情况。
			// 此处可以选择返回已存在文件的信息，并考虑是否删除刚刚通过预签名URL上传的重复文件。
			l.Logger.Infof("确认上传时发现文件哈希 %s 已存在 (ID: %d)，用户ID: %d。新上传对象键: %s", in.Hash, existingFile.Id, in.UserId, in.FileKey)

			if existingFile.FilePath != in.FileKey { // 确保不是同一个对象
				go func() { // 异步删除，不阻塞当前请求
					if delErr := l.svcCtx.Minio.DeleteFile(context.Background(), in.FileKey); delErr != nil {
						l.Logger.Errorf("删除预签名上传的重复对象 %s 失败: %v", in.FileKey, delErr)
					} else {
						l.Logger.Infof("成功删除预签名上传的重复对象: %s", in.FileKey)
					}
				}()
			}

			return &file.ConfirmUploadResponse{
				FileId:  existingFile.Id,
				FileUrl: existingFile.FileUrl,
				Success: true,
				Message: "文件已存在，确认成功 (秒传)",
			}, nil
		}
	}

	// 4. 生成文件访问 URL
	fileURL := l.buildMinioFileUrl(in.FileKey)

	// 5. 保存文件记录到数据库
	fileRecord := model.Files{
		Filename:     filepath.Base(in.FileKey), // MinIO 中的实际文件名 (通常是 pkgminio.GenerateObjectName 生成的)
		OriginalName: in.Filename,               // 用户上传的原始文件名
		FilePath:     in.FileKey,                // MinIO Object Key
		FileUrl:      fileURL,
		FileType:     in.FileType,
		FileSize:     in.FileSize, // 使用从MinIO获取或客户端提供的大小
		MimeType:     in.MimeType, // 使用从MinIO获取或客户端提供的MIME
		Hash:         in.Hash,     // 客户端计算并提供的哈希
		UserId:       in.UserId,
		Status:       1,
		CreateAt:     time.Now().Unix(),
		UpdateAt:     time.Now().Unix(),
	}

	if err := l.svcCtx.DB.Create(&fileRecord).Error; err != nil {
		l.Logger.Errorf("保存文件记录到数据库失败 (ConfirmUpload): %v, 文件信息: %+v", err, fileRecord)
		// 数据库保存失败，MinIO中已存在文件。根据策略，可以尝试删除，或标记为待处理。
		// 暂时不删除 MinIO 文件，因为可能只是DB暂时故障。
		return &file.ConfirmUploadResponse{
			Success: false,
			Message: "保存文件信息失败: " + err.Error(),
		}, nil
	}

	cacheKey := fmt.Sprintf("presigned:%s", in.FileKey)
	l.svcCtx.Redis.Del(l.ctx, cacheKey)

	return &file.ConfirmUploadResponse{
		FileId:  fileRecord.Id,
		FileUrl: fileRecord.FileUrl,
		Success: true,
		Message: "确认上传成功",
	}, nil
}

// 构建 MinIO 文件 URL 的辅助函数 (可以考虑放到公共包或 base logic)
func (l *ConfirmUploadLogic) buildMinioFileUrl(objectName string) string {
	cfgMinio := l.svcCtx.Config.MinIO
	cfgStorage := l.svcCtx.Config.FileStorage
	if cfgStorage.BaseURL != "" {
		return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(cfgStorage.BaseURL, "/"), cfgMinio.BucketName, objectName)
	}
	scheme := "http"
	if cfgMinio.UseSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, cfgMinio.Endpoint, cfgMinio.BucketName, objectName)
}
