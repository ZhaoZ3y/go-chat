package logic

import (
	"context"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取文件列表
func (l *GetFileListLogic) GetFileList(in *file.GetFileListRequest) (*file.GetFileListResponse, error) {
	// todo: add your logic here and delete this line

	return &file.GetFileListResponse{}, nil
}
