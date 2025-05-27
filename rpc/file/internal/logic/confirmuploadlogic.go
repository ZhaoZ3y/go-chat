package logic

import (
	"context"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmUploadLogic {
	return &ConfirmUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 确认上传
func (l *ConfirmUploadLogic) ConfirmUpload(in *file.ConfirmUploadRequest) (*file.ConfirmUploadResponse, error) {
	// todo: add your logic here and delete this line

	return &file.ConfirmUploadResponse{}, nil
}
