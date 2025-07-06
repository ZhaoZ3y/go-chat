package logic

import (
	"IM/pkg/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileRecordLogic {
	return &GetFileRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户文件记录
func (l *GetFileRecordLogic) GetFileRecord(in *file.GetFileRecordReq) (*file.GetFileRecordResp, error) {
	// 查询用户的文件记录
	var fileRecords []model.FileRecord
	err := l.svcCtx.DB.WithContext(l.ctx).Where("user_id = ?", in.UserId).Find(&fileRecords).Error
	if err != nil {
		l.Logger.Errorf("查询用户文件记录失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询文件记录失败")
	}

	// 如果没有找到文件记录，返回空列表
	if len(fileRecords) == 0 {
		return &file.GetFileRecordResp{
			FileRecords: []*file.GetFileInfoResponse{},
		}, nil
	}

	var fileInfos []*file.GetFileInfoResponse
	var currentTime int64
	if val := l.ctx.Value("currentTime"); val != nil {
		if ts, ok := val.(int64); ok {
			currentTime = ts
		}
	}

	for _, record := range fileRecords {
		isExpired := currentTime != 0 && record.ExpireAt < currentTime

		fileInfos = append(fileInfos, &file.GetFileInfoResponse{
			FileId:      record.FileID,
			FileName:    record.FileName,
			FileSize:    record.FileSize,
			ContentType: record.ContentType,
			UserId:      record.UserID,
			CreatedAt:   record.CreateAt,
			ExpireAt:    record.ExpireAt,
			IsExpired:   isExpired,
			Etag:        record.ETag,
			FileType:    record.FileType,
		})
	}

	return &file.GetFileRecordResp{
		FileRecords: fileInfos,
	}, nil
}
