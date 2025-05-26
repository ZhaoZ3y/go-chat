package logic

import (
	"context"

	"IM/apps/file/internal/svc"
	"IM/apps/file/rpc/file"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageLogic {
	return &UploadImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传图片
func (l *UploadImageLogic) UploadImage(in *file.UploadImageRequest) (*file.UploadImageResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadImageResponse{}, nil
}
