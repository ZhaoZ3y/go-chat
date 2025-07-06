package logic

import (
	"IM/pkg/model"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"io"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 下载文件
func (l *DownloadFileLogic) DownloadFile(in *file.DownloadFileRequest) (*file.DownloadFileResponse, error) {
	// 1. 查找数据库记录
	var fileRecord model.FileRecord
	result := l.svcCtx.DB.WithContext(l.ctx).Where("file_id = ?", in.FileId).First(&fileRecord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "文件未找到")
		}
		l.Logger.Errorf("从数据库查询文件信息失败: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "数据库错误")
	}

	// 2. 从MinIO下载文件对象
	// 警告：直接在服务中流式传输大文件会消耗大量内存和带宽
	l.Logger.Infof("正在为文件 %s (大小: %d) 提供直接下载服务", fileRecord.ObjectName, fileRecord.FileSize)
	object, err := l.svcCtx.MinioClient.DownloadFile(l.ctx, fileRecord.ObjectName)
	if err != nil {
		l.Logger.Errorf("从MinIO下载对象 %s 失败: %v", fileRecord.ObjectName, err)
		return nil, status.Errorf(codes.Internal, "下载文件失败")
	}
	defer object.Close()

	// 3. 将文件内容读取到字节切片中
	fileData, err := io.ReadAll(object)
	if err != nil {
		l.Logger.Errorf("读取对象 %s 的数据流失败: %v", fileRecord.ObjectName, err)
		return nil, status.Errorf(codes.Internal, "读取文件数据失败")
	}

	// 4. 返回响应
	return &file.DownloadFileResponse{
		FileData:    fileData,
		FileName:    fileRecord.FileName,
		FileSize:    fileRecord.FileSize,
		ContentType: fileRecord.ContentType,
	}, nil
}
