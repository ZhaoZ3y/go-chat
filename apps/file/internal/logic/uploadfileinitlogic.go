package logic

import (
	"context"

	"IM/apps/file/internal/svc"
	"IM/apps/file/rpc/file"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileInitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileInitLogic {
	return &UploadFileInitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传文件初始化
func (l *UploadFileInitLogic) UploadFileInit(in *file.UploadFileInitRequest) (*file.UploadFileInitResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadFileInitResponse{}, nil
}
