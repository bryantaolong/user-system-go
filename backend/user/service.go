package user

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/bryantaolong/user-system/auth"
	"github.com/bryantaolong/user-system/types"
	"github.com/bryantaolong/user-system/pkg/jwt"
	"github.com/bryantaolong/user-system/response"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// RoleService 角色服务（依赖 auth 模块的角色仓库）
type RoleService struct {
	roleRepo *auth.UserRoleRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo *auth.UserRoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

// ListAll 查询所有角色
func (s *RoleService) ListAll(ctx context.Context) ([]model.UserRole, error) {
	return s.roleRepo.SelectAll()
}

// GetDefaultRole 查询默认角色
func (s *RoleService) GetDefaultRole(ctx context.Context) (*model.UserRole, error) {
	return s.roleRepo.SelectOneByIsDefaultTrue()
}

// ListByIDs 根据 ID 列表查询角色
func (s *RoleService) ListByIDs(ctx context.Context, ids []int) ([]model.UserRole, error) {
	return s.roleRepo.SelectByIDList(ids)
}

// Service 用户服务
type Service struct {
	userRepo *auth.UserRepository
	roleSvc  *RoleService
	logger   *zap.Logger
}

// NewService 创建用户服务实例
func NewService(
	userRepo *auth.UserRepository,
	roleSvc *RoleService,
	logger *zap.Logger,
) *Service {
	return &Service{
		userRepo: userRepo,
		roleSvc:  roleSvc,
		logger:   logger,
	}
}

// CreateUser 管理员创建用户
func (s *Service) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.SysUser, error) {
	existing, _ := s.userRepo.SelectByUsername(req.Username)
	if existing != nil {
		return nil, response.NewBusinessError("用户名已存在")
	}

	roleIDs := req.RoleIDs
	if len(roleIDs) == 0 {
		defaultRole, err := s.roleSvc.GetDefaultRole(ctx)
		if err != nil || defaultRole == nil {
			return nil, response.NewBusinessError("系统未配置默认角色")
		}
		roleIDs = []int{defaultRole.ID}
	}

	roles, err := s.roleSvc.ListByIDs(ctx, roleIDs)
	if err != nil {
		return nil, response.NewBusinessError("角色查询失败")
	}
	if len(roles) != len(roleIDs) {
		return nil, response.NewResourceNotFoundError("部分角色不存在")
	}

	roleNames := make([]string, 0, len(roles))
	for _, r := range roles {
		roleNames = append(roleNames, r.RoleName)
	}
	rolesStr := strings.Join(roleNames, ",")

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessError("密码加密失败")
	}

	now := time.Now()
	operator := getCurrentOperator(ctx)
	user := &model.SysUser{
		Username:       req.Username,
		Password:       string(hashedPwd),
		Phone:          req.Phone,
		Email:          req.Email,
		Roles:          &rolesStr,
		Status:         model.UserStatusNormal,
		LoginFailCount: intPtr(0),
		Deleted:        0,
		Version:        0,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      &operator,
		UpdatedBy:      &operator,
	}

	if err := s.userRepo.Insert(user); err != nil {
		return nil, response.NewPersistenceError("插入数据库失败")
	}

	s.logger.Info("管理员创建用户成功", zap.Int64("id", user.ID), zap.String("username", user.Username))
	return user, nil
}

// GetAllUsers 获取所有用户列表（分页）
func (s *Service) GetAllUsers(ctx context.Context, pageNum, pageSize int) (*response.PageResponse, error) {
	offset := (pageNum - 1) * pageSize
	users, err := s.userRepo.SelectPageByConditions(offset, pageSize, nil)
	if err != nil {
		return nil, response.NewBusinessError("查询用户列表失败")
	}
	total, _ := s.userRepo.Count(nil)

	rows := make([]any, 0, len(users))
	for i := range users {
		rows = append(rows, users[i])
	}
	pr := response.NewPageResponse(rows, total, pageNum, pageSize)
	return &pr, nil
}

// GetUserByID 根据用户 ID 获取用户信息
func (s *Service) GetUserByID(ctx context.Context, userID int64) (*model.SysUser, error) {
	user, err := s.userRepo.SelectByID(userID)
	if err != nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}
	return user, nil
}

// GetUserByUsername 根据用户名获取用户信息
func (s *Service) GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	user, err := s.userRepo.SelectByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}
	return user, nil
}

// QueryUsers 通用用户搜索
func (s *Service) QueryUsers(ctx context.Context, req *model.UserQueryRequest, pageNum, pageSize int) (*response.PageResponse, error) {
	offset := (pageNum - 1) * pageSize
	users, err := s.userRepo.SelectPageByConditions(offset, pageSize, req)
	if err != nil {
		return nil, response.NewBusinessError("查询用户失败")
	}
	total, _ := s.userRepo.Count(req)

	rows := make([]any, 0, len(users))
	for i := range users {
		rows = append(rows, users[i])
	}
	pr := response.NewPageResponse(rows, total, pageNum, pageSize)
	return &pr, nil
}

// UpdateUser 更新用户基础信息
func (s *Service) UpdateUser(ctx context.Context, userID int64, req *model.UserUpdateRequest) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Email != nil {
		user.Email = req.Email
	}

	user.Version++
	user.UpdatedAt = time.Now()
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewPersistenceError("更新失败，可能数据已变更")
	}

	s.logger.Info("用户信息更新成功", zap.Int64("userId", userID))
	return user, nil
}

// ChangeRoleByIds 修改用户角色
func (s *Service) ChangeRoleByIds(ctx context.Context, userID int64, req *model.ChangeRoleRequest) (*model.SysUser, error) {
	roles, err := s.roleSvc.ListByIDs(ctx, req.RoleIDs)
	if err != nil {
		return nil, response.NewBusinessError("角色查询失败")
	}
	if len(roles) != len(req.RoleIDs) {
		return nil, response.NewResourceNotFoundError("部分角色不存在")
	}

	roleNames := make([]string, 0, len(roles))
	for _, r := range roles {
		roleNames = append(roleNames, r.RoleName)
	}
	rolesStr := strings.Join(roleNames, ",")

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Roles = &rolesStr
	user.Version++
	user.UpdatedAt = time.Now()
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewPersistenceError("角色更新失败")
	}
	return user, nil
}

// ResetPassword 重置用户密码（管理员）
func (s *Service) ResetPassword(ctx context.Context, userID int64, newPassword string) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessError("密码加密失败")
	}

	now := time.Now()
	user.Password = string(hashedPwd)
	user.PasswordResetAt = &now
	user.Version++
	user.UpdatedAt = now
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewPersistenceError("密码重置失败")
	}

	s.logger.Info("用户密码强制修改成功", zap.Int64("userId", userID))
	return user, nil
}

// BlockUser 封禁指定用户
func (s *Service) BlockUser(ctx context.Context, userID int64) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}

	user.Status = model.UserStatusBanned
	user.Version++
	user.UpdatedAt = time.Now()
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewPersistenceError("封禁失败")
	}

	s.logger.Info("用户封禁成功", zap.Int64("userId", userID))
	return user, nil
}

// UnblockUser 解封指定用户
func (s *Service) UnblockUser(ctx context.Context, userID int64) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}

	user.Status = model.UserStatusNormal
	user.Version++
	user.UpdatedAt = time.Now()
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewPersistenceError("解封失败")
	}

	s.logger.Info("用户解封成功", zap.Int64("userId", userID))
	return user, nil
}

// DeleteUser 删除用户（逻辑删除）
func (s *Service) DeleteUser(ctx context.Context, userID int64) (int64, error) {
	claims, _ := ctx.Value("claims").(*jwt.Claims)
	operator := "0"
	if claims != nil && claims.Username != "" {
		operator = claims.Username
	}
	if err := s.userRepo.DeleteByID(userID, operator); err != nil {
		return 0, response.NewResourceNotFoundError("用户不存在或已被删除")
	}
	s.logger.Info("用户删除成功 (逻辑删除)", zap.Int64("userId", userID))
	return userID, nil
}

// ExportAllUsers 全量导出用户数据为 Excel
func (s *Service) ExportAllUsers(ctx context.Context, pageSize int) (*excelize.File, error) {
	if pageSize <= 0 {
		pageSize = 1000
	}

	f := excelize.NewFile()
	sheet := "用户列表"
	f.SetSheetName("Sheet1", sheet)

	// 设置表头
	headers := []string{"ID", "用户名", "手机", "邮箱", "状态", "角色", "最后登录时间", "最后登录IP", "最后登录设备", "密码重置时间", "删除状态", "创建时间", "更新时间", "创建人", "更新人"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 分批查询写入
	pageNum := 1
	row := 2

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		offset := (pageNum - 1) * pageSize
		users, err := s.userRepo.SelectPageByConditions(offset, pageSize, nil)
		if err != nil {
			return nil, err
		}
		if len(users) == 0 {
			break
		}

		for _, user := range users {
			vo := model.ToExportVO(&user)
			cols := []any{
				vo.ID, vo.Username, derefStr(vo.Phone), derefStr(vo.Email),
				vo.Status, derefStr(vo.Roles),
				formatTime(vo.LastLoginAt), derefStr(vo.LastLoginIP),
				derefStr(vo.LastLoginDevice), formatTime(vo.PasswordResetAt),
				vo.Deleted, vo.CreatedAt.Format("2006-01-02 15:04:05"),
				vo.UpdatedAt.Format("2006-01-02 15:04:05"),
				derefStr(vo.CreatedBy), derefStr(vo.UpdatedBy),
			}
			for i, v := range cols {
				cell, _ := excelize.CoordinatesToCellName(i+1, row)
				f.SetCellValue(sheet, cell, v)
			}
			row++
		}

		s.logger.Info("已导出数据", zap.Int("count", row-2))
		pageNum++
	}

	// 设置表头样式
	style, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#CCFFCC"}},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheet, "A1", "O1", style)

	s.logger.Info("导出完成", zap.Int("total", row-2))
	return f, nil
}

// ProfileService 用户资料服务
type ProfileService struct {
	profileRepo *ProfileRepository
	fileSvc     *FileService
	logger      *zap.Logger
}

// NewProfileService 创建用户资料服务实例
func NewProfileService(
	profileRepo *ProfileRepository,
	fileSvc *FileService,
	logger *zap.Logger,
) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		fileSvc:     fileSvc,
		logger:      logger,
	}
}

// CreateUserProfile 创建用户资料
func (s *ProfileService) CreateUserProfile(ctx context.Context, record *model.UserProfile) (*model.UserProfile, error) {
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
func (s *ProfileService) GetUserProfileByUserId(ctx context.Context, userID int64) (*model.UserProfile, error) {
	profile, err := s.profileRepo.SelectByUserID(userID)
	if err != nil {
		return nil, response.NewResourceNotFoundError("用户信息不存在")
	}
	return profile, nil
}

// GetUserProfileByUserIdOrEmpty 如果不存在则返回空实体
func (s *ProfileService) GetUserProfileByUserIdOrEmpty(ctx context.Context, userID int64) *model.UserProfile {
	profile, err := s.profileRepo.SelectByUserID(userID)
	if err != nil {
		s.logger.Warn("用户资料不存在，返回空实体", zap.Int64("userId", userID))
		return &model.UserProfile{UserID: userID}
	}
	return profile
}

// GetUserProfileByRealName 根据真实姓名查询用户资料
func (s *ProfileService) GetUserProfileByRealName(ctx context.Context, realName string) (*model.UserProfile, error) {
	profile, err := s.profileRepo.SelectByRealName(realName)
	if err != nil {
		return nil, response.NewResourceNotFoundError("用户信息不存在")
	}
	return profile, nil
}

// UpdateUserProfile 更新用户资料
func (s *ProfileService) UpdateUserProfile(ctx context.Context, userID int64, req *model.UserUpdateRequest) (*model.UserProfile, error) {
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
		t, err := time.Parse(time.RFC3339, *req.Birthday)
		if err != nil {
			return nil, response.NewBusinessError("出生日期格式不正确")
		}
		if t.After(time.Now()) {
			return nil, response.NewBusinessError("出生日期必须是过去的日期")
		}
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
func (s *ProfileService) UpdateAvatar(ctx context.Context, userID int64, file *multipart.FileHeader) (string, error) {
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
		if !s.fileSvc.DeleteFile(*profile.Avatar) {
			s.logger.Warn("删除旧头像失败", zap.String("avatar", *profile.Avatar), zap.Int64("userId", userID))
		}
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

// getCurrentOperator 获取当前操作者标识（用户ID）
func getCurrentOperator(ctx context.Context) string {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return "0"
	}
	return fmt.Sprintf("%d", claims.UserID)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }
