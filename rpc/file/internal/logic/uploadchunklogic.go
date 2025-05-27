package logic

import (
	"context"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadChunkLogic {
	return &UploadChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分块上传
func (l *UploadChunkLogic) UploadChunk(in *file.UploadChunkRequest) (*file.UploadChunkResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadChunkResponse{}, nil
}
