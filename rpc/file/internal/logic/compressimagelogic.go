package logic

import (
	pkgminio "IM/pkg/minio"
	"IM/pkg/model"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"
	"time"

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
	var originalFile model.Files
	if err := l.svcCtx.DB.Where("id = ? AND status = 1", in.FileId).First(&originalFile).Error; err != nil {
		return &file.CompressImageResponse{
			Success: false,
			Message: "原始文件不存在或已被删除",
		}, nil
	}

	// 2. 检查是否为图片文件
	if originalFile.FileType != "image" {
		return &file.CompressImageResponse{
			Success: false,
			Message: "指定的文件不是图片类型",
		}, nil
	}

	// 3. 从MinIO下载原始图片
	minioObject, err := l.svcCtx.Minio.DownloadFile(l.ctx, originalFile.FilePath)
	if err != nil {
		l.Logger.Errorf("从 MinIO 下载原始图片 %s 失败: %v", originalFile.FilePath, err)
		return &file.CompressImageResponse{
			Success: false,
			Message: "读取原始图片失败: " + err.Error(),
		}, nil
	}
	defer minioObject.Close()

	originalImageData, err := io.ReadAll(minioObject) // 读取到内存，因为 image.Decode 需要 io.Reader, resize 也可能需要
	if err != nil {
		l.Logger.Errorf("读取 MinIO 对象 %s 数据失败: %v", originalFile.FilePath, err)
		return &file.CompressImageResponse{Success: false, Message: "读取图片数据失败"}, nil
	}
	minioObject.Close() // 明确关闭

	// 4. 解码图片
	img, format, err := image.Decode(bytes.NewReader(originalImageData))
	if err != nil {
		l.Logger.Errorf("解码图片 %s (原始格式: %s) 失败: %v", originalFile.OriginalName, originalFile.MimeType, err)
		return &file.CompressImageResponse{
			Success: false,
			Message: "图片解码失败，可能格式不支持: " + err.Error(),
		}, nil
	}

	// 5. 调整图片尺寸
	var resizedImg image.Image = img
	targetWidth := uint(in.Width)
	targetHeight := uint(in.Height)

	if targetWidth > 0 || targetHeight > 0 {
		// 如果一个维度为0，则根据另一个维度等比缩放
		if targetWidth == 0 && targetHeight > 0 {
			originalWidth := uint(img.Bounds().Dx())
			originalHeight := uint(img.Bounds().Dy())
			if originalHeight > 0 {
				targetWidth = uint(float64(originalWidth) * (float64(targetHeight) / float64(originalHeight)))
			} else {
				targetWidth = originalWidth // 容错
			}
		} else if targetHeight == 0 && targetWidth > 0 {
			originalWidth := uint(img.Bounds().Dx())
			originalHeight := uint(img.Bounds().Dy())
			if originalWidth > 0 {
				targetHeight = uint(float64(originalHeight) * (float64(targetWidth) / float64(originalWidth)))
			} else {
				targetHeight = originalHeight // 容错
			}
		}
		resizedImg = resize.Resize(targetWidth, targetHeight, img, resize.Lanczos3)
	}

	// 6. 压缩图片到缓冲区
	var buf bytes.Buffer
	quality := int(in.Quality)
	if quality <= 0 || quality > 100 {
		quality = 75 // 默认压缩质量
	}

	outputMimeType := originalFile.MimeType // 默认使用原始MIME类型
	switch strings.ToLower(format) {        // format 是从 image.Decode 获取的实际格式
	case "jpeg", "jpg":
		options := &jpeg.Options{Quality: quality}
		err = jpeg.Encode(&buf, resizedImg, options)
		outputMimeType = "image/jpeg"
	case "png":
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression} // 可调整 DefaultCompression
		err = encoder.Encode(&buf, resizedImg)
		outputMimeType = "image/png"
	case "gif": // GIF 压缩比较复杂，标准库不支持有损压缩或优化
		err = gif.Encode(&buf, resizedImg, nil)
		outputMimeType = "image/gif"
	default:
		// 对于其他格式，尝试统一输出为 JPEG
		l.Logger.Errorf("图片 %s 的原始格式 %s 不直接支持质量压缩，将尝试转为JPEG输出", originalFile.OriginalName, format)
		options := &jpeg.Options{Quality: quality}
		err = jpeg.Encode(&buf, resizedImg, options)
		outputMimeType = "image/jpeg"
	}

	if err != nil {
		l.Logger.Errorf("图片编码/压缩 %s 失败: %v", originalFile.OriginalName, err)
		return &file.CompressImageResponse{
			Success: false,
			Message: "图片编码/压缩失败: " + err.Error(),
		}, nil
	}
	compressedData := buf.Bytes()

	// 7. 生成压缩后的文件名和MinIO对象名
	// 构造一个能体现压缩参数的文件名
	ext := filepath.Ext(originalFile.OriginalName)
	if outputMimeType == "image/jpeg" && (ext != ".jpg" && ext != ".jpeg") {
		ext = ".jpg"
	} else if outputMimeType == "image/png" && ext != ".png" {
		ext = ".png"
	}
	baseOriginalName := strings.TrimSuffix(originalFile.OriginalName, filepath.Ext(originalFile.OriginalName))
	compressedOriginalName := fmt.Sprintf("%s_compressed_q%d_w%d_h%d%s", baseOriginalName, quality, targetWidth, targetHeight, ext)
	compressedObjectKey := pkgminio.GenerateObjectName(originalFile.UserId, "image_compressed", compressedOriginalName)

	// 8. 上传压缩后的图片到MinIO
	err = l.svcCtx.Minio.UploadFile(l.ctx, compressedObjectKey, bytes.NewReader(compressedData), int64(len(compressedData)), outputMimeType)
	if err != nil {
		l.Logger.Errorf("上传压缩后的图片 %s 到 MinIO 失败: %v", compressedObjectKey, err)
		return &file.CompressImageResponse{
			Success: false,
			Message: "上传压缩文件失败: " + err.Error(),
		}, nil
	}

	// 9. 生成压缩文件URL
	compressedUrl := l.buildMinioFileUrl(compressedObjectKey)

	// 10. 计算压缩后文件的哈希
	compressedHash := fmt.Sprintf("%x", md5.Sum(compressedData))

	// 11. 保存压缩文件信息到数据库 (作为一条新的文件记录)
	compressedRecord := model.Files{
		Filename:     filepath.Base(compressedObjectKey), // MinIO中的实际文件名
		OriginalName: compressedOriginalName,             // 描述性的原始名
		FilePath:     compressedObjectKey,                // MinIO Object Key
		FileUrl:      compressedUrl,
		FileType:     "image", // 压缩后仍是图片
		FileSize:     int64(len(compressedData)),
		MimeType:     outputMimeType,
		Hash:         compressedHash, // 压缩后内容的哈希
		UserId:       originalFile.UserId,
		Status:       1,
		CreateAt:     time.Now().Unix(),
		UpdateAt:     time.Now().Unix(),
	}

	if err := l.svcCtx.DB.Create(&compressedRecord).Error; err != nil {
		l.Logger.Errorf("保存压缩文件记录 %s 到数据库失败: %v", compressedObjectKey, err)
		// 回滚：删除已上传到 MinIO 的压缩文件
		if delErr := l.svcCtx.Minio.DeleteFile(l.ctx, compressedObjectKey); delErr != nil {
			l.Logger.Errorf("数据库保存压缩记录失败后，删除 MinIO 对象 %s 也失败: %v", compressedObjectKey, delErr)
		}
		return &file.CompressImageResponse{
			Success: false,
			Message: "保存压缩文件信息失败: " + err.Error(),
		}, nil
	}

	return &file.CompressImageResponse{
		CompressedFileId:  compressedRecord.Id,
		CompressedFileUrl: compressedUrl,
		Success:           true,
		Message:           "图片压缩成功",
	}, nil
}

// 构建 MinIO 文件 URL 的辅助函数
func (l *CompressImageLogic) buildMinioFileUrl(objectName string) string {
	cfgMinio := l.svcCtx.Config.MinIO
	cfgStorage := l.svcCtx.Config.FileStorage
	if cfgStorage.BaseURL != "" {
		return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(cfgStorage.BaseURL, "/"), cfgMinio.BucketName, objectName)
	}
	scheme := "http"
	if cfgMinio.UseSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, cfgMinio.Endpoint, cfgMinio.BucketName, objectName)
}
