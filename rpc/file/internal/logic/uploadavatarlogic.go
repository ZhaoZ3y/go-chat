package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/fileutil"
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传头像
func (l *UploadAvatarLogic) UploadAvatar(in *file.UploadFileRequest) (*file.UploadFileResponse, error) {
	fileID := uuid.NewString()
	reader := bytes.NewReader(in.FileData)

	// 构建头像文件名
	objectName := fmt.Sprintf("avatar/%d_%d_%s", in.UserId, time.Now().Unix(), in.FileName)
	fileType := fileutil.GetFileTypeFromName(in.FileName)

	// 上传到 MinIO（avatar 前缀）
	uploadInfo, err := l.svcCtx.MinioClient.UploadFile(
		l.ctx,
		objectName,
		reader,
		in.FileSize,
		in.ContentType,
	)
	if err != nil {
		return nil, err
	}

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
		ExpireAt:    0,
	}

	// 写入数据库
	result := l.svcCtx.DB.WithContext(l.ctx).Create(fileRecord)
	if result.Error != nil {
		l.Logger.Errorf("在数据库中创建文件记录失败: %v", result.Error)
		// 尽力清理已上传的孤儿文件
		_ = l.svcCtx.MinioClient.DeleteFile(context.Background(), objectName)
		return nil, status.Errorf(codes.Internal, "保存文件元数据失败")
	}

	// 构造公开访问地址（MinIO 需配置桶为 public-read）
	publicURL := fmt.Sprintf("http://%s/%s/%s",
		l.svcCtx.Config.MinIO.Endpoint,
		l.svcCtx.Config.MinIO.BucketName,
		objectName,
	)

	// 更新用户头像字段
	err = l.svcCtx.DB.Model(&model.User{}).
		Where("id = ?", in.UserId).
		Update("avatar", publicURL).Error
	if err != nil {
		l.Errorf("更新用户头像失败: %v", err)
		return nil, err
	}

	return &file.UploadFileResponse{
		FileId:   objectName,
		FileName: in.FileName,
		FileSize: in.FileSize,
		ExpireAt: 0, // 头像不自动过期
		FileUrl:  publicURL,
	}, nil
}
