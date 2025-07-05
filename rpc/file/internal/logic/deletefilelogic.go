package logic

import (
	"IM/pkg/model"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除文件
func (l *DeleteFileLogic) DeleteFile(in *file.DeleteFileRequest) (*file.DeleteFileResponse, error) {
	// 1. 先查找文件记录，以便进行权限校验
	var fileRecord model.FileRecord
	result := l.svcCtx.DB.WithContext(l.ctx).Where("file_id = ?", in.FileId).First(&fileRecord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 文件不存在，也认为删除成功，保证幂等性
			return &file.DeleteFileResponse{Success: true}, nil
		}
		l.Logger.Errorf("查找待删除文件记录失败: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "数据库错误")
	}

	// 2. 权限校验：确保是文件所有者在执行删除操作
	if fileRecord.UserID != in.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "无权删除此文件")
	}

	// 3. 执行软删除
	delResult := l.svcCtx.DB.WithContext(l.ctx).Where("file_id = ?", in.FileId).Delete(&model.FileRecord{})
	if delResult.Error != nil {
		l.Logger.Errorf("从数据库软删除文件记录失败: %v", delResult.Error)
		return nil, status.Errorf(codes.Internal, "删除文件失败")
	}

	l.Logger.Infof("已软删除文件记录，ID: %s", in.FileId)
	return &file.DeleteFileResponse{Success: true}, nil
}
