package logic

import (
	"context"

	"IM/apps/file/internal/svc"
	"IM/apps/file/rpc/file"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传头像
func (l *UploadAvatarLogic) UploadAvatar(in *file.UploadAvatarRequest) (*file.UploadAvatarResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadAvatarResponse{}, nil
}
