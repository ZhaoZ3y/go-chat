package logic

import (
	pkgminio "IM/pkg/minio"
	build "IM/pkg/utils/file"
	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GeneratePresignedUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGeneratePresignedUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneratePresignedUrlLogic {
	return &GeneratePresignedUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成预签名URL
func (l *GeneratePresignedUrlLogic) GeneratePresignedUrl(in *file.GeneratePresignedUrlRequest) (*file.GeneratePresignedUrlResponse, error) {
	// 1. 确定文件类型
	fileType := build.GetFileTypeFromName(in.Filename)

	objectName := pkgminio.GenerateObjectName(in.UserId, fileType, in.Filename)

	// 3. 设置过期时间
	expires := time.Duration(in.ExpireTime) * time.Second
	if expires <= 0 {
		expires = 15 * time.Minute // 默认15分钟过期
	}
	if expires > 7*24*time.Hour { // MinIO 预签名URL最长7天
		expires = 7 * 24 * time.Hour
	}

	// 4. 生成预签名上传 URL
	uploadURL, err := l.svcCtx.Minio.GetPresignedPutURL(l.ctx, objectName, expires)
	if err != nil {
		l.Logger.Errorf("为对象 %s 生成预签名上传URL失败: %v", objectName, err)
		return &file.GeneratePresignedUrlResponse{
			Success: false,
			Message: "生成上传凭证失败: " + err.Error(),
		}, nil
	}

	return &file.GeneratePresignedUrlResponse{
		UploadUrl:  uploadURL,
		FileKey:    objectName, // 客户端上传时需要使用此 Key
		ExpireTime: time.Now().Add(expires).Unix(),
		Success:    true,
		Message:    "生成成功",
	}, nil
}
