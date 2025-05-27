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

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户注册
func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 验证输入参数
	if in.Username == "" || in.Password == "" || in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "用户名、密码和邮箱不能为空")
	}

	// 检查用户名是否已存在
	var existUser model.User
	err := l.svcCtx.DB.Where("username = ? AND deleted_at = 0", in.Username).First(&existUser).Error
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "用户名已存在")
	}
	if err != gorm.ErrRecordNotFound {
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 检查邮箱是否已存在
	err = l.svcCtx.DB.Where("email = ? AND deleted_at = 0", in.Email).First(&existUser).Error
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "邮箱已存在")
	}
	if err != gorm.ErrRecordNotFound {
		l.Logger.Errorf("查询邮箱失败: %v", err)
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 加密密码
	password := encrypt.EncryptPassword(in.Password)

	// 创建用户
	newUser := model.User{
		Username: in.Username,
		Password: password,
		Nickname: in.Nickname,
		Email:    in.Email,
		Status:   1,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}

	if newUser.Nickname == "" {
		newUser.Nickname = in.Username
	}

	err = l.svcCtx.DB.Create(&newUser).Error
	if err != nil {
		l.Logger.Errorf("创建用户失败: %v", err)
		return nil, status.Error(codes.Internal, "注册失败")
	}

	return &user.RegisterResponse{
		Message: "注册成功",
	}, nil
}
