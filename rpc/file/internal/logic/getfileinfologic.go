package logic

import (
	"IM/pkg/model"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"

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
	var fileRecord model.FileRecord
	result := l.svcCtx.DB.WithContext(l.ctx).Where("file_id = ?", in.FileId).First(&fileRecord)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "文件未找到")
		}
		l.Logger.Errorf("从数据库查询文件信息失败: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "数据库错误")
	}

	isExpired := time.Now().Unix() > fileRecord.ExpireAt

	// 组装响应
	return &file.GetFileInfoResponse{
		FileId:      fileRecord.FileID,
		FileName:    fileRecord.FileName,
		FileSize:    fileRecord.FileSize,
		ContentType: fileRecord.ContentType,
		UserId:      fileRecord.UserID,
		CreatedAt:   fileRecord.CreateAt,
		ExpireAt:    fileRecord.ExpireAt,
		IsExpired:   isExpired,
		Etag:        fileRecord.ETag,
		FileType:    fileRecord.FileType,
	}, nil
}
