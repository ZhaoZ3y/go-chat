package logic

import (
	"context"

	"IM/rpc/file/file"
	"IM/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileInfoLogic {
	return &GetFileInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取文件信息
func (l *GetFileInfoLogic) GetFileInfo(in *file.GetFileInfoRequest) (*file.GetFileInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &file.GetFileInfoResponse{}, nil
}
