package logic

import (
	"context"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

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
	// todo: add your logic here and delete this line

	return &file.GeneratePresignedUrlResponse{}, nil
}
