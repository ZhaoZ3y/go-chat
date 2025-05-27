package logic

import (
	"context"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompressImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompressImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompressImageLogic {
	return &CompressImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 图片压缩
func (l *CompressImageLogic) CompressImage(in *file.CompressImageRequest) (*file.CompressImageResponse, error) {
	// todo: add your logic here and delete this line

	return &file.CompressImageResponse{}, nil
}
