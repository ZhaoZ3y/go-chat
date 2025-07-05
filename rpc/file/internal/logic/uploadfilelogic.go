package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/fileutil"
	"bytes"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

const fileExpireDays = 7

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

// 上传文件
func (l *UploadFileLogic) UploadFile(in *file.UploadFileRequest) (*file.UploadFileResponse, error) {
	fileID := uuid.NewString()
	objectName := fileID // 使用UUID作为对象名保证唯一性
	reader := bytes.NewReader(in.FileData)

	// 上传文件到MinIO并获取ETag
	uploadInfo, err := l.svcCtx.MinioClient.UploadFile(l.ctx, objectName, reader, in.FileSize, in.ContentType)
	if err != nil {
		l.Logger.Errorf("上传文件到MinIO失败: %v", err)
		return nil, status.Errorf(codes.Internal, "上传文件失败")
	}

	// 获取文件类型
	fileType := fileutil.GetFileTypeFromName(in.FileName)

	// 准备数据库记录
	expireAt := time.Now().Add(fileExpireDays * 24 * time.Hour).Unix()
	fileRecord := &model.FileRecord{
		FileID:      fileID,
		FileName:    in.FileName,
		FileType:    fileType,
		FileSize:    in.FileSize,
		ContentType: in.ContentType,
		ETag:        uploadInfo.ETag,
		ObjectName:  objectName,
		UserID:      in.UserId,
		Status:      1,
		ExpireAt:    expireAt,
	}

	// 写入数据库
	result := l.svcCtx.DB.WithContext(l.ctx).Create(fileRecord)
	if result.Error != nil {
		l.Logger.Errorf("在数据库中创建文件记录失败: %v", result.Error)
		// 尽力清理已上传的孤儿文件
		_ = l.svcCtx.MinioClient.DeleteFile(context.Background(), objectName)
		return nil, status.Errorf(codes.Internal, "保存文件元数据失败")
	}

	// 返回符合proto定义的响应
	return &file.UploadFileResponse{
		FileId:   fileRecord.FileID,
		FileName: fileRecord.FileName,
		FileSize: fileRecord.FileSize,
		ExpireAt: fileRecord.ExpireAt,
	}, nil
}
