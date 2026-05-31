package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bryan/user-system/internal/config"
	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/repository"
	"go.uber.org/zap"
)

type UserProfileService struct {
	profileRepo *repository.UserProfileRepository
	fileSvc     *LocalFileService
	logger      *zap.Logger
}

func NewUserProfileService(
	profileRepo *repository.UserProfileRepository,
	fileSvc *LocalFileService,
	logger *zap.Logger,
) *UserProfileService {
	return &UserProfileService{
		profileRepo: profileRepo,
		fileSvc:     fileSvc,
		logger:      logger,
	}
}

// CreateUserProfile 创建用户资料
func (s *UserProfileService) CreateUserProfile(ctx context.Context, record *model.UserProfile) (*model.UserProfile, error) {
	now := time.Now()
	operator := getCurrentOperator(ctx)
	record.Deleted = 0
	record.Version = 0
	record.CreatedAt = now
	record.UpdatedAt = now
	record.CreatedBy = &operator
	record.UpdatedBy = &operator

	if err := s.profileRepo.Insert(record); err != nil {
		return nil, response.NewPersistenceError("创建用户信息失败")
	}
	s.logger.Info("用户信息创建成功", zap.Int64("userId", record.UserID))
	return record, nil
}

// GetUserProfileByUserId 根据用户主键查询用户资料
func (s *UserProfileService) GetUserProfileByUserId(ctx context.Context, userID int64) (*model.UserProfile, error) {
	profile, err := s.profileRepo.SelectByUserID(userID)
	if err != nil {
		return nil, response.NewResourceNotFoundError("用户信息不存在")
	}
	return profile, nil
}

// GetUserProfileByUserIdOrEmpty 如果不存在则返回空实体
func (s *UserProfileService) GetUserProfileByUserIdOrEmpty(ctx context.Context, userID int64) *model.UserProfile {
	profile, err := s.profileRepo.SelectByUserID(userID)
	if err != nil {
		s.logger.Warn("用户资料不存在，返回空实体", zap.Int64("userId", userID))
		return &model.UserProfile{UserID: userID}
	}
	return profile
}

// GetUserProfileByRealName 根据真实姓名查询用户资料
func (s *UserProfileService) GetUserProfileByRealName(ctx context.Context, realName string) (*model.UserProfile, error) {
	profile, err := s.profileRepo.SelectByRealName(realName)
	if err != nil {
		return nil, response.NewResourceNotFoundError("用户信息不存在")
	}
	return profile, nil
}

// UpdateUserProfile 更新用户资料
func (s *UserProfileService) UpdateUserProfile(ctx context.Context, userID int64, req *model.UserUpdateRequest) (*model.UserProfile, error) {
	profile, err := s.GetUserProfileByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.RealName != nil {
		profile.RealName = req.RealName
	}
	if req.Gender != nil {
		g := model.Gender(*req.Gender)
		profile.Gender = &g
	}
	if req.Birthday != nil {
		t, _ := time.Parse(time.RFC3339, *req.Birthday)
		profile.Birthday = &t
	}
	if req.Avatar != nil {
		profile.Avatar = req.Avatar
	}

	profile.Version++
	profile.UpdatedAt = time.Now()
	profile.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.profileRepo.Update(profile); err != nil {
		return nil, response.NewPersistenceError("用户信息更新失败")
	}

	s.logger.Info("用户信息更新成功", zap.Int64("userId", userID))
	return profile, nil
}

// UpdateAvatar 上传并更新用户头像
func (s *UserProfileService) UpdateAvatar(ctx context.Context, userID int64, file *multipart.FileHeader) (string, error) {
	profile, err := s.GetUserProfileByUserId(ctx, userID)
	if err != nil {
		return "", err
	}

	avatarPath, err := s.fileSvc.StoreFile(file, "avatars")
	if err != nil {
		return "", response.NewPersistenceError("头像上传失败: " + err.Error())
	}

	// 删除旧头像
	if profile.Avatar != nil && *profile.Avatar != "" {
		s.fileSvc.DeleteFile(*profile.Avatar)
	}

	profile.Avatar = &avatarPath
	profile.Version++
	profile.UpdatedAt = time.Now()
	profile.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.profileRepo.Update(profile); err != nil {
		return "", response.NewPersistenceError("头像更新失败")
	}

	s.logger.Info("用户头像更新成功", zap.Int64("userId", userID), zap.String("path", avatarPath))
	return avatarPath, nil
}

// LocalFileService 文件存储服务
type LocalFileService struct {
	uploadDir string
	logger    *zap.Logger
}

func NewLocalFileService(logger *zap.Logger) *LocalFileService {
	dir := "./uploads"
	if config.AppConfig != nil && config.AppConfig.File.UploadDir != "" {
		dir = config.AppConfig.File.UploadDir
	}
	return &LocalFileService{uploadDir: dir, logger: logger}
}

func (s *LocalFileService) StoreFile(file *multipart.FileHeader, subDir string) (string, error) {
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

func (s *LocalFileService) DeleteFile(filePath string) bool {
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
	return false
}
