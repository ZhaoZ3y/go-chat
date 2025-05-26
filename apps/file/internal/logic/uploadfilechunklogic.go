package logic

import (
	"context"

	"IM/apps/file/internal/svc"
	"IM/apps/file/rpc/file"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileChunkLogic {
	return &UploadFileChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传文件分片
func (l *UploadFileChunkLogic) UploadFileChunk(in *file.UploadFileChunkRequest) (*file.UploadFileChunkResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadFileChunkResponse{}, nil
}
