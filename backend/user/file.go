package user

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bryantaolong/user-system/config"
	"github.com/bryantaolong/user-system/response"
	"go.uber.org/zap"
)

// FileService 本地文件存储服务
type FileService struct {
	uploadDir string
	logger    *zap.Logger
}

// NewFileService 创建文件存储服务实例
func NewFileService(logger *zap.Logger) *FileService {
	dir := "./uploads"
	if config.AppConfig != nil && config.AppConfig.File.UploadDir != "" {
		dir = config.AppConfig.File.UploadDir
	}
	return &FileService{uploadDir: dir, logger: logger}
}

// StoreFile 存储文件到本地
func (s *FileService) StoreFile(file *multipart.FileHeader, subDir string) (string, error) {
	uploadPath := filepath.Join(s.uploadDir, subDir)
	absPath, _ := filepath.Abs(uploadPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return "", err
		}
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 验证文件内容类型
	buf := make([]byte, 8)
	n, _ := src.Read(buf)
	if n < 2 {
		return "", io.ErrUnexpectedEOF
	}
	if !isValidImageType(buf[:n]) {
		return "", response.NewBusinessError("不支持的文件类型，仅允许 PNG、JPEG、GIF、WebP 格式")
	}
	src.Seek(0, 0)

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".png"
	}
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixMilli(), strings.ReplaceAll(file.Filename, ext, "")+ext)
	dst := filepath.Join(absPath, fileName)

	dstFile, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, src); err != nil {
		return "", err
	}

	return filepath.Join(subDir, fileName), nil
}

// DeleteFile 删除本地文件
func (s *FileService) DeleteFile(filePath string) bool {
	fullPath := filepath.Join(s.uploadDir, filePath)
	err := os.Remove(fullPath)
	if err != nil {
		s.logger.Error("删除文件失败", zap.String("path", filePath), zap.Error(err))
		return false
	}
	return true
}

func isValidImageType(header []byte) bool {
	// PNG: 89 50 4E 47 0D 0A 1A 0A
	if len(header) >= 8 && header[0] == 0x89 && header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47 {
		return true
	}
	// JPEG: FF D8
	if len(header) >= 2 && header[0] == 0xFF && header[1] == 0xD8 {
		return true
	}
	// GIF: 47 49 46 38 (GIF8)
	if len(header) >= 4 && header[0] == 0x47 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x38 {
		return true
	}
	// WebP: 52 49 46 46 ... 57 45 42 50 (RIFF...WEBP)
	if len(header) >= 12 && header[0] == 0x52 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x46 && header[8] == 0x57 && header[9] == 0x45 && header[10] == 0x42 && header[11] == 0x50 {
		return true
	}
	return false
}
