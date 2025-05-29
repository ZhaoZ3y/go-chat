package logic

import (
	"IM/pkg/model"
	"context"
	"gorm.io/gorm"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileInfoLogic {
	return &GetFileInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取文件信息
func (l *GetFileInfoLogic) GetFileInfo(in *file.GetFileInfoRequest) (*file.GetFileInfoResponse, error) {
	var fileRecord model.Files
	query := l.svcCtx.DB.Where("id = ? AND status = ?", in.FileId, 1) // Status 1 表示正常
	if in.UserId > 0 {                                                // 假设 UserId > 0 时表示需要用户权限校验
		query = query.Where("user_id = ?", in.UserId)
	}

	if err := query.First(&fileRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &file.GetFileInfoResponse{
				Success: false,
				Message: "文件不存在、已被删除或无权查看",
			}, nil
		}
		l.Logger.Errorf("查询文件信息失败: %v, FileID: %d, UserID: %d", err, in.FileId, in.UserId)
		return &file.GetFileInfoResponse{
			Success: false,
			Message: "查询文件信息时发生错误",
		}, nil
	}

	return &file.GetFileInfoResponse{
		FileInfo: &file.FileInfo{
			Id:           fileRecord.Id,
			Filename:     fileRecord.Filename,     // MinIO中的实际文件名
			OriginalName: fileRecord.OriginalName, // 用户上传的原始文件名
			FilePath:     fileRecord.FilePath,     // MinIO Object Key
			FileUrl:      fileRecord.FileUrl,      // 访问URL
			FileType:     fileRecord.FileType,
			FileSize:     fileRecord.FileSize,
			MimeType:     fileRecord.MimeType,
			Hash:         fileRecord.Hash,
			UserId:       fileRecord.UserId,
			Status:       int32(fileRecord.Status),
			CreateAt:     fileRecord.CreateAt,
			UpdateAt:     fileRecord.UpdateAt,
		},
		Success: true,
		Message: "获取成功",
	}, nil
}
