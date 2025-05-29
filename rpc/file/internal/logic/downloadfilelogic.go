package logic

import (
	"IM/pkg/model"
	"context"
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
	// 1. 获取文件记录
	var fileRecord model.Files
	query := l.svcCtx.DB.Where("id = ? AND status = ?", in.FileId, 1) // Status 1 表示正常
	if in.UserId > 0 {                                                // 假设 UserId > 0 时表示需要用户权限校验
		query = query.Where("user_id = ?", in.UserId)
	}

	if err := query.First(&fileRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &file.DownloadFileResponse{
				Success: false,
				Message: "文件不存在、已被删除或无权访问",
			}, nil
		}
		l.Logger.Errorf("查询待下载文件记录失败: %v, FileID: %d, UserID: %d", err, in.FileId, in.UserId)
		return &file.DownloadFileResponse{
			Success: false,
			Message: "查询文件信息失败",
		}, nil
	}

	object, err := l.svcCtx.Minio.DownloadFile(l.ctx, fileRecord.FilePath)
	if err != nil {
		l.Logger.Errorf("从 MinIO 下载文件 %s 失败: %v", fileRecord.FilePath, err)
		return &file.DownloadFileResponse{
			Success: false,
			Message: "下载文件失败: " + err.Error(),
		}, nil
	}
	defer object.Close()

	fileData, err := io.ReadAll(object)
	if err != nil {
		l.Logger.Errorf("读取 MinIO 对象 %s 数据流失败: %v", fileRecord.FilePath, err)
		return &file.DownloadFileResponse{
			Success: false,
			Message: "读取文件数据失败: " + err.Error(),
		}, nil
	}

	return &file.DownloadFileResponse{
		FileData: fileData,
		Filename: fileRecord.OriginalName, // 返回原始文件名，方便客户端保存
		MimeType: fileRecord.MimeType,
		Success:  true,
		Message:  "下载成功",
	}, nil
}
