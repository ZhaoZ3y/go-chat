package logic

import (
	"IM/pkg/model"
	"IM/pkg/utils/encrypt"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"

	"IM/rpc/user/internal/svc"
	"IM/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改密码
func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	if in.UserId == 0 || in.OldPassword == "" || in.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "参数不能为空")
	}
	if in.OldPassword == in.NewPassword {
		return nil, status.Error(codes.InvalidArgument, "新旧密码不能相同")
	}

	var userModel model.User
	err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", in.UserId).First(&userModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 校验旧密码
	encryptedOld := encrypt.EncryptPassword(in.OldPassword)
	if userModel.Password != encryptedOld {
		return nil, status.Error(codes.Unauthenticated, "旧密码错误")
	}

	// 校验用户状态
	if userModel.Status != 1 {
		return nil, status.Error(codes.PermissionDenied, "用户已被禁用")
	}

	// 更新密码
	newEncrypted := encrypt.EncryptPassword(in.NewPassword)
	err = l.svcCtx.DB.Model(&userModel).Updates(map[string]interface{}{
		"password":  newEncrypted,
		"update_at": time.Now().Unix(),
	}).Error
	if err != nil {
		l.Logger.Errorf("更新密码失败: %v", err)
		return nil, status.Error(codes.Internal, "更新密码失败")
	}

	return &user.ChangePasswordResponse{
		Message: "密码修改成功",
	}, nil
}
