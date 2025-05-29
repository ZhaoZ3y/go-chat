package logic

import (
	"IM/pkg/model"
	"context"
	"gorm.io/gorm"
	"time"

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
	// 1. 获取文件记录
	var fileRecord model.Files
	// 确保检查 UserId 以进行权限控制
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ? AND status = ?", in.FileId, in.UserId, 1).First(&fileRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &file.DeleteFileResponse{
				Success: false,
				Message: "文件不存在或已被删除，或无权操作",
			}, nil
		}
		l.Logger.Errorf("查询待删除文件记录失败: %v, FileID: %d, UserID: %d", err, in.FileId, in.UserId)
		return &file.DeleteFileResponse{
			Success: false,
			Message: "查询文件信息失败",
		}, nil
	}

	// 2. 软删除数据库记录
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		l.Logger.Errorf("开启数据库事务失败 (DeleteFile): %v", tx.Error)
		return &file.DeleteFileResponse{Success: false, Message: "删除操作失败"}, nil
	}

	// 更新 Status 和 DeletedAt
	updateData := map[string]interface{}{
		"status":     2, // 标记为已删除
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"update_at":  time.Now().Unix(),
	}
	if err := tx.Model(&model.Files{}).Where("id = ?", fileRecord.Id).Updates(updateData).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("软删除数据库文件记录失败: %v, FileID: %d", err, fileRecord.Id)
		return &file.DeleteFileResponse{
			Success: false,
			Message: "更新文件状态失败",
		}, nil
	}

	// 3. 从MinIO删除文件对象
	if err := l.svcCtx.Minio.DeleteFile(l.ctx, fileRecord.FilePath); err != nil {
		tx.Rollback()
		l.Logger.Errorf("从 MinIO 删除文件对象 %s 失败: %v。数据库操作已回滚。", fileRecord.FilePath, err)
		return &file.DeleteFileResponse{
			Success: false,
			Message: "从存储服务删除文件失败",
		}, nil
	}

	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交数据库事务失败 (DeleteFile): %v", err)
		return &file.DeleteFileResponse{Success: false, Message: "删除操作最终确认失败"}, nil
	}

	return &file.DeleteFileResponse{
		Success: true,
		Message: "删除成功",
	}, nil
}
