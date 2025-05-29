package logic

import (
	"IM/pkg/model"
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
	query := l.svcCtx.DB.Model(&model.Files{}).Where("user_id = ? AND status = ?", in.UserId, 1)

	// 文件类型过滤
	if in.FileType != "" {
		query = query.Where("file_type = ?", in.FileType)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Logger.Errorf("查询文件列表总数失败: %v, UserID: %d, FileType: %s", err, in.UserId, in.FileType)
		return &file.GetFileListResponse{Total: 0, Files: []*file.FileInfo{}}, nil // 返回空列表和0总数
	}

	if total == 0 {
		return &file.GetFileListResponse{Total: 0, Files: []*file.FileInfo{}}, nil
	}

	// 分页参数处理
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10 // 默认每页10条
	}
	page := in.Page
	if page <= 0 {
		page = 1 // 默认第一页
	}
	offset := (page - 1) * pageSize

	var files []model.Files
	// 查询分页数据，并按创建时间降序排列
	if err := query.Order("create_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&files).Error; err != nil {
		l.Logger.Errorf("查询文件列表分页数据失败: %v, UserID: %d, FileType: %s, Page: %d, PageSize: %d", err, in.UserId, in.FileType, page, pageSize)
		return &file.GetFileListResponse{Total: total, Files: []*file.FileInfo{}}, nil // 返回总数，但文件列表为空
	}

	// 将查询结果转换为响应格式
	var fileInfos []*file.FileInfo
	for _, f := range files {
		fileInfos = append(fileInfos, &file.FileInfo{
			Id:           f.Id,
			Filename:     f.Filename,
			OriginalName: f.OriginalName,
			FilePath:     f.FilePath,
			FileUrl:      f.FileUrl,
			FileType:     f.FileType,
			FileSize:     f.FileSize,
			MimeType:     f.MimeType,
			Hash:         f.Hash,
			UserId:       f.UserId,
			Status:       int32(f.Status),
			CreateAt:     f.CreateAt,
			UpdateAt:     f.UpdateAt,
		})
	}

	return &file.GetFileListResponse{
		Files: fileInfos,
		Total: total,
	}, nil
}
