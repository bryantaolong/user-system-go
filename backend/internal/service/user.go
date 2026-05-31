package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/pkg/jwt"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/repository"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleSvc  *UserRoleService
	logger   *zap.Logger
}

func NewUserService(
	userRepo *repository.UserRepository,
	roleSvc *UserRoleService,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleSvc:  roleSvc,
		logger:   logger,
	}
}

// CreateUser 管理员创建用户
func (s *UserService) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.SysUser, error) {
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
func (s *UserService) GetAllUsers(ctx context.Context, pageNum, pageSize int) (*response.PageResult, error) {
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
	pr := response.NewPageResult(rows, total, pageNum, pageSize)
	return &pr, nil
}

// GetUserByID 根据用户 ID 获取用户信息
func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*model.SysUser, error) {
	user, err := s.userRepo.SelectByID(userID)
	if err != nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}
	return user, nil
}

// GetUserByUsername 根据用户名获取用户信息
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	user, err := s.userRepo.SelectByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUsers 通用用户搜索
func (s *UserService) QueryUsers(ctx context.Context, req *model.UserQueryRequest, pageNum, pageSize int) (*response.PageResult, error) {
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
	pr := response.NewPageResult(rows, total, pageNum, pageSize)
	return &pr, nil
}

// UpdateUser 更新用户基础信息
func (s *UserService) UpdateUser(ctx context.Context, userID int64, req *model.UserUpdateRequest) (*model.SysUser, error) {
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
func (s *UserService) ChangeRoleByIds(ctx context.Context, userID int64, req *model.ChangeRoleRequest) (*model.SysUser, error) {
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
func (s *UserService) ResetPassword(ctx context.Context, userID int64, newPassword string) (*model.SysUser, error) {
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
func (s *UserService) BlockUser(ctx context.Context, userID int64) (*model.SysUser, error) {
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
func (s *UserService) UnblockUser(ctx context.Context, userID int64) (*model.SysUser, error) {
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
func (s *UserService) DeleteUser(ctx context.Context, userID int64) (int64, error) {
	operator := getCurrentOperator(ctx)
	if err := s.userRepo.DeleteByID(userID, operator); err != nil {
		return 0, response.NewResourceNotFoundError("用户不存在或已被删除")
	}
	s.logger.Info("用户删除成功 (逻辑删除)", zap.Int64("userId", userID))
	return userID, nil
}

// ExportAllUsers 全量导出用户数据为 Excel
func (s *UserService) ExportAllUsers() (*excelize.File, error) {
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
	pageSize := 1000
	row := 2

	for {
		offset := (pageNum - 1) * pageSize
		users, err := s.userRepo.SelectPageByConditions(offset, pageSize, nil)
		if err != nil || len(users) == 0 {
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

// getCurrentOperator 获取当前操作者标识
func getCurrentOperator(ctx context.Context) string {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return "SYSTEM"
	}
	return fmt.Sprintf("%s", claims.Username)
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
