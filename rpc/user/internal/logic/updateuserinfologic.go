package logic

import (
	"IM/pkg/model"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoReq) (*user.UpdateUserInfoResp, error) {
	updates := make(map[string]interface{})

	if in.Nickname != "" {
		updates["nickname"] = in.Nickname
	}
	if in.Avatar != "" {
		updates["avatar"] = in.Avatar
	}
	if in.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var existUser model.User
		err := l.svcCtx.DB.Where("email = ? AND id != ?", in.Email, in.UserId).First(&existUser).Error
		if err == nil {
			return nil, fmt.Errorf("邮箱已被使用")
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		updates["email"] = in.Email
	}

	if len(updates) > 0 {
		updates["update_at"] = time.Now()
		err := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", in.UserId).Updates(updates).Error
		if err != nil {
			return nil, err
		}
	}

	return &user.UpdateUserInfoResp{}, nil
}
