package file

import (
	"path/filepath"
	"strings"
)

// GetFileTypeFromName 根据文件名获取文件类型
func GetFileTypeFromName(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return "image"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm":
		return "video"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a":
		return "audio"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md":
		return "document"
	case ".zip", ".rar", ".tar", ".gz", ".7z":
		return "archive"
	default:
		return "other"
	}
}
